// +build gofuzz

package fuzz

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pragmatwice/go-squirrel/comparer"

	// to pin dep in go.mod
	_ "github.com/oraluben/go-fuzz/go-fuzz-dep"

	"github.com/pingcap/tidb/config"
	"github.com/pingcap/tidb/tidb-server/internal"
)

var tidbConn *sql.DB = nil
var mysqlConn *sql.DB = nil

var err error

func init() {
	os.Args = []string{os.Args[0]}

	instanceDir, err := ioutil.TempDir("", "tidb-fuzz.*")
	if err != nil {
		panic(err)
	}
	sockName := path.Join(instanceDir, "tidb.sock")
	storeDir := path.Join(instanceDir, "store")
	tmpDir := path.Join(instanceDir, "tmp")
	logFile := path.Join(instanceDir, "tidb.log")

	go internal.MainWithConfig(func(c *config.Config) {
		c.Host = ""
		c.Port = 0
		c.Socket = sockName
		c.Store = "unistore"
		c.Path = storeDir
		c.Status.StatusPort = 0
		c.TempStoragePath = tmpDir
		c.Log.File.Filename = logFile
	})

	mysqlInstanceDir := strings.ReplaceAll(instanceDir, "tidb-fuzz", "mysql-fuzz")
	err = os.Mkdir(mysqlInstanceDir, os.ModePerm)
	if err != nil {
		panic(err)
	}

	mysqlSockName := path.Join(mysqlInstanceDir, "mysql.sock")
	mysqlPidFile := path.Join(mysqlInstanceDir, "mysql.pid")
	mysqlDataDir := path.Join(mysqlInstanceDir, "data")

	// ref to https://dev.mysql.com/doc/refman/8.0/en/multiple-servers.html
	mysqldInit := exec.Command("mysqld", "--initialize-insecure", fmt.Sprintf("--datadir=%s", mysqlDataDir))
	err = mysqldInit.Run()
	if err != nil {
		panic(err)
	}

	mysqld := exec.Command("mysqld",
		fmt.Sprintf("--datadir=%s", mysqlDataDir),
		fmt.Sprintf("--pid-file=%s", mysqlPidFile),
		fmt.Sprintf("--socket=%s", mysqlSockName),
		"--port=0")

	err = mysqld.Start()
	if err != nil {
		panic(err)
	}

	tc, mc := make(chan *sql.DB), make(chan *sql.DB)

	go sqlConnect(sockName, tc)
	go sqlConnect(mysqlSockName, mc)

	tidbConn, mysqlConn = <-tc, <-mc

	fmt.Println(instanceDir)
}

func sqlConnect(sockName string, cc chan *sql.DB) {
	var conn *sql.DB

	for i := 0; i < 5; i++ {
		if _, err := os.Stat(sockName); os.IsNotExist(err) {
			time.Sleep(time.Second)
			continue
		}

		conn, err = sql.Open("mysql", fmt.Sprintf("root@unix(%s)/", sockName))
		if err != nil {
			time.Sleep(time.Second)
		} else {
			break
		}
	}

	if err != nil {
		panic(fmt.Sprintf("%s not up after 5 seconds", sockName))
	}

	_, _ = conn.Exec("create database test") // useful in mysql

	_, err = conn.Exec("use test")
	if err != nil {
		panic(err)
	}

	cc <- conn
}

func isSelect(sql string) bool {
	sql = strings.TrimLeft(sql, " (\n")
	sql = strings.ToLower(sql)
	return strings.HasPrefix(sql, "select") || strings.HasPrefix(sql, "with")
}

// Fuzz is the required name by go-fuzz
func Fuzz(raw []byte) int {
	query := string(raw)

	// fmt.Println(query)
	tidbErr, mysqlErr := make(chan error), make(chan error)

	if isSelect(query) {
		exec := func(conn *sql.DB, rows **sql.Rows, ec chan error) {
			var err error
			*rows, err = conn.Query(query)
			ec <- err
		}

		var tidbRows, mysqlRows *sql.Rows

		go exec(tidbConn, &tidbRows, tidbErr)
		go exec(mysqlConn, &mysqlRows, mysqlErr)

		te := <-tidbErr
		me := <-mysqlErr

		if te != nil || me != nil {
			if te != nil && me != nil {
				panic(fmt.Sprintf("[both err] tidb: %v; mysql: %v", te, me))
			} else {
				panic(fmt.Sprintf("[one side err] tidb: %v; mysql: %v", te, me))
			}
		}

		tidbSR, err := comparer.NewSortedRows(tidbRows)
		if err != nil {
			panic(err)
		}

		mysqlSR, err := comparer.NewSortedRows(mysqlRows)
		if err != nil {
			panic(err)
		}

		k, l, r := comparer.DiffMetaInfo(tidbSR, mysqlSR)
		if k != comparer.NoDiff {
			panic(fmt.Sprintf("[metainfo diff (%v)] tidb: %v, mysql: %v\n", k, l, r))
		}

		dr := comparer.DiffData(tidbSR, mysqlSR)
		if dr != nil {
			panic(fmt.Sprintf("[data diff] %v", dr))
		}

	} else {
		exec := func(conn *sql.DB, ec chan error) {
			_, err := conn.Exec(query)
			ec <- err
		}

		go exec(tidbConn, tidbErr)
		go exec(mysqlConn, mysqlErr)

		te := <-tidbErr
		me := <-mysqlErr

		// assume that ddls are correct
		if te != nil || me != nil {
			panic(fmt.Sprintf("[ddl error] tidb: %v; mysql: %v", te, me))
		}
	}

	return 1
}

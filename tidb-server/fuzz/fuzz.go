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

	fmt.Println(instanceDir)

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
	mysqlDataDir := path.Join(mysqlInstanceDir, "data")

	// ref to https://dev.mysql.com/doc/refman/8.0/en/multiple-servers.html
	mysqldInit := exec.Command("mysqld", "--initialize", fmt.Sprintf("--datadir=%s", mysqlDataDir))
	err = mysqldInit.Run()
	if err != nil {
		panic(err)
	}

	mysqld := exec.Command("mysqld", fmt.Sprintf("--datadir=%s", mysqlDataDir), fmt.Sprintf("--socket=%s", mysqlSockName), "--port=0")
	err = mysqld.Start()
	if err != nil {
		panic(err)
	}

	tc := make(chan *sql.DB)
	mc := make(chan *sql.DB)

	go sqlConnect(sockName, tc)
	go sqlConnect(mysqlSockName, mc)

	tidbConn, mysqlConn = <-tc, <-mc
}

func sqlConnect(sockName string, cc chan *sql.DB) {
	var conn *sql.DB

	for i := 0; i < 5; i++ {
		conn, err = sql.Open("mysql", fmt.Sprintf("root@unix(%s)/test", sockName))
		if err != nil {
			time.Sleep(time.Second)
			continue
		}
	}

	if err != nil {
		panic(fmt.Sprintf("%s not up after 5 seconds", sockName))
	}

	cc <- conn
}

func isSelect(sql string) bool {
	sql = strings.TrimLeft(sql, " (\n")
	sql = strings.ToLower(sql)
	return strings.HasPrefix(sql, "select")
}

// Fuzz is the required name by go-fuzz
func Fuzz(raw []byte) int {
	query := string(raw)

	// fmt.Println(query)
	tidbErr, mysqlErr := make(chan error), make(chan error)

	if isSelect(query) {
		exec := func(conn *sql.DB, rc chan *sql.Rows, ec chan error) {
			rows, err := conn.Query(query)
			rc <- rows
			ec <- err
		}

		tidbRows, mysqlRows := make(chan *sql.Rows), make(chan *sql.Rows)

		go exec(tidbConn, tidbRows, tidbErr)
		go exec(mysqlConn, mysqlRows, mysqlErr)

		err = <-tidbErr
		if err != nil {
			fmt.Println(err)
			return 0
		}

		err = <-mysqlErr
		if <-mysqlErr != nil {
			fmt.Println(err)
			return 0
		}

		tidbSR, err := comparer.NewSortedRows(<-tidbRows)
		if err != nil {
			fmt.Println(err)
			return 0
		}

		mysqlSR, err := comparer.NewSortedRows(<-mysqlRows)
		if err != nil {
			fmt.Println(err)
			return 0
		}

		k, l, r := comparer.DiffMetaInfo(tidbSR, mysqlSR)
		if k != comparer.NoDiff {
			fmt.Printf("[metainfo diff (%v)] tidb: %v, mysql: %v\n", k, l, r)
			return 0
		}

		dr := comparer.DiffData(tidbSR, mysqlSR)
		if dr != nil {
			fmt.Printf("[data diff] %v", dr)
			return 0
		}

	} else {
		exec := func(conn *sql.DB, ec chan error) {
			_, err := conn.Exec(query)
			ec <- err
		}

		go exec(tidbConn, tidbErr)
		go exec(mysqlConn, mysqlErr)

		err = <-tidbErr
		if err != nil {
			fmt.Println(err)
			return 0
		}

		err = <-mysqlErr
		if err != nil {
			fmt.Println(err)
			return 0
		}
	}

	return 1
}

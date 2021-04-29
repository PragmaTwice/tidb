// +build gofuzz

package fuzz

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"time"

	_ "github.com/go-sql-driver/mysql"
	// to pin dep in go.mod
	_ "github.com/oraluben/go-fuzz/go-fuzz-dep"

	"github.com/pingcap/tidb/config"
	"github.com/pingcap/tidb/tidb-server/internal"
)

var conn *sql.DB = nil
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

	go internal.MainWithConfig(func(c *config.Config) {
		c.Host = ""
		c.Port = 0
		c.Socket = sockName
		c.Store = "unistore"
		c.Path = storeDir
		c.Status.StatusPort = 0
		c.TempStoragePath = tmpDir
	})

	for i := 0; i < 5; i++ {
		conn, err = sql.Open("mysql", fmt.Sprintf("root@unix(%s)/test", sockName))
		if err != nil {
			time.Sleep(time.Second)
			continue
		}
	}
	if err != nil {
		panic("TiDB not up after 5 seconds")
	}
}

// Fuzz is the required name by go-fuzz
func Fuzz(raw []byte) int {
	query := string(raw)

	_, err = conn.Exec(query)

	if err != nil {
		return 0
	}

	return 1
}

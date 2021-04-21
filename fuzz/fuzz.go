// +build gofuzz

package fuzz

import server "github.com/pingcap/tidb/tidb-server"

func init() {
	server.Main()
}

// Fuzz is the required name by go-fuzz
func Fuzz(_ []byte) int {
	return 0
}

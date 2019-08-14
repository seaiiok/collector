package fileparse

import (
	"fmt"
	"gcom/gtools/gdb"
	"testing"
	"time"

	_ "github.com/alexbrainman/odbc"
)

func TestFile(t *testing.T) {

	dbserver := "192.168.1.18"
	dbuser := "sa"
	dbpw := "kmtSoft12345678"
	dbase := "ifixsvr"
	dbport := "1433"

	dsn := fmt.Sprintf("driver={sql server};server=%s;port=%s;uid=%s;pwd=%s;database=%s;encrypt=disable", dbserver, dbport, dbuser, dbpw, dbase)

	err := gdb.New("odbc", dsn)
	if err != nil {
		t.Log("odbc init err:", err)
	}

	gdb.Exec(sql_createtable3)

	time.Sleep(time.Second)
	err = Run()

	t.Log(err)

}

package appserver

import (
	"collector/pkg/interfaces"
	"fmt"
	"gcom/gtools/gdb"
	"sync"

	_ "github.com/alexbrainman/odbc"
)

var once sync.Once

type app struct {
	config interfaces.IConfig
	log    interfaces.ILog
}

func New(config interfaces.IConfig, log interfaces.ILog) *app {
	return &app{
		config: config,
		log:    log,
	}
}

func (this *app) Run() {
	NewServer(this.config, this.log)
	once.Do(func() {
		dsn := fmt.Sprintf("driver={sql server};server=%s;port=%s;uid=%s;pwd=%s;database=%s;encrypt=disable", this.config.Get("dbserver"), this.config.Get("dbport"), this.config.Get("dbuser"), this.config.Get("dbpw"), this.config.Get("dbase"))
		this.log.Info("odbc driver:", dsn)
		err := gdb.New("odbc", dsn)
		if err != nil {
			this.log.Error("odbc init err:", err)
		}
	})

	gdb.Exec(sql_createtable1)
	gdb.Exec(sql_createtable2)
	gdb.Exec(sql_createtable3)

}

package appserver

import (
	"collector/internal/pkg/apptray"
	"collector/internal/pkg/fileparse"
	"collector/pkg/interfaces"
	"context"
	"fmt"
	"gcom/gtools/gdb"
	"os"
	"sync"

	_ "github.com/alexbrainman/odbc"
)

var once sync.Once

const (
	readPool = "./read-pool/"
)

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

func (this *app) Run(ctx context.Context, cancel context.CancelFunc) {
	NewServer(ctx, this.config, this.log)

	once.Do(func() {
		os.MkdirAll(readPool, 0744)
		dsn := fmt.Sprintf("driver={sql server};server=%s;port=%s;uid=%s;pwd=%s;database=%s;encrypt=disable", this.config.Get("dbserver"), this.config.Get("dbport"), this.config.Get("dbuser"), this.config.Get("dbpw"), this.config.Get("dbase"))
		err := gdb.New("odbc", dsn)
		if err != nil {
			this.log.Error("odbc init err:", err)
		}
	})

	gdb.Exec(sql_createtable1)
	gdb.Exec(sql_createtable2)
	gdb.Exec(sql_createtable3)

	this.log.Info("Collect Online!")

	fileparse.Run(ctx, readPool)

	at := apptray.New(ctx, cancel, this.config)
	at.AppTray()
}

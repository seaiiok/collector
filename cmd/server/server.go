package main

import (
	"fmt"
	"os"
	"time"

	"github.com/seaiiok/gcom"
	"github.com/seaiiok/snet/snet.v1"
)

var (
	g          = gcom.New()
	configFile = "./logcollect.json"
	Config     = make(map[string]interface{}, 0)
)

func init() {
	var err error
	Config, err = g.GConfig.Config2Map(configFile)
	if err != nil {
		os.Exit(0)
	}
}
func main() {

	go func() {
		svr := snet.New(Config["host"].(string), Config["port"].(string))

		svr.OnConnect(func(conn *snet.Connection) {

		})

		svr.OnRecvMessage(func(conn *snet.Connection, msg []byte) {
			fmt.Println(msg)
			conn.OnSendMsg(msg)
		})

		svr.OnSendMessage(func(conn *snet.Connection) {

		})

		svr.OnDisConnect(func(conn *snet.Connection) {

		})

	}()

	for {
		time.Sleep(10 * time.Second)
	}
}

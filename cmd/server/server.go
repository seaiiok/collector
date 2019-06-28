package main

import (
	"fmt"
	"time"

	"github.com/seaiiok/snet/snet.v1"
)

func init() {

}
func main() {

	go func() {
		svr := snet.New("127.0.0.1", "495")

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

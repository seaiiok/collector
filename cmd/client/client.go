package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/seaiiok/gcom"
	"github.com/seaiiok/snet/snet.v1"
)

var (
	g          = gcom.New()
	configFile = "./logcollect.json"
	Config     = make(map[string]interface{}, 0)
	netWork    = "tcp4"
)

func init() {
	var err error
	Config, err = g.GConfig.Config2Map(configFile)
	if err != nil {
		os.Exit(0)
	}
}

func main() {
	conn, err := net.Dial(netWork, Config["host"].(string)+":"+Config["port"].(string))
	if err != nil {
		fmt.Println("client start err, exit!")
		return
	}

	go ReadBuf(conn)

	for {
		time.Sleep(1 * time.Second)
	}
}

func ReadBuf(conn net.Conn) {
	p := snet.Package{}
	for i := 0; i < 10; i++ {
		// fmt.Println("发送1：",MakeSomeMsg(i))
		// fmt.Println("发送1：",p.Pack(MakeSomeMsg(i)))
		_, err := conn.Write(p.Pack(MakeSomeMsg(i)))
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("Recv goroutine running...")
		buf := make([]byte, 1024)
		cnt, err := conn.Read(buf)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("接收完整包：", buf[:cnt])
		b1 := p.UnPack(buf[:cnt])
		fmt.Println("接收数据：", p.UnPack(buf[:cnt]))
		fmt.Println("接收：", string(b1))

		// time.Sleep(1 * time.Second)
	}

}

func MakeSomeMsg(i int) []byte {
	x := "hello world " + strconv.Itoa(i)
	return []byte(x)
}

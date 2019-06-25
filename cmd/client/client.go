package main

import (
	"fmt"
	"net"

	"github.com/seaiiok/snet/snet.v1"
)

func main() {
	clientGo()
}

//建一个客户端
func clientGo() {
	conn, err := net.Dial("tcp", "127.0.0.1:495")
	if err != nil {
		fmt.Println("client dial err, exit!")
		return
	}
for i := 0; i < 100; i++ {
	//封包
	msg1 := makeSomeMsg(byte(i))
	b := msg1.Pack()

	_, err = conn.Write(b)
	if err != nil {
		fmt.Println(err)
	}

	buf := make([]byte, 1024)
	cnt, err := conn.Read(buf)
	if err != nil {
		fmt.Println(err)
		return
	}

	//解包
	msg2 := snet.Package{}
	msg := msg2.UnPack(buf[:cnt])
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("------------接收数据------------")
	fmt.Println("ID:", msg.ID)
	fmt.Println("Key长度:", msg.KeyLen, "Key内容:", msg.Key)
	fmt.Println("Data长度:", msg.DataLen, "Data内容:", msg.Data)
}

}

//制造一些数据
func makeSomeMsg(id byte) snet.Package {
	return snet.Package{
		ID:      id,
		KeyLen:  0,
		DataLen: 0,
		Key:     []string{"golang", "tcp"},
		Data:    [][]string{{"1", "data1"}, {"2", "data2"}, {"2", "data2"}},
	}
}

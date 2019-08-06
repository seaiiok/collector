package netclient

import (
	"collector/api"
	"collector/pkg/network/netclient"
	"gcom/gcmd"
	"net"

	"github.com/golang/protobuf/proto"
)

type client struct {
}

func NewClient(ip string, port int) *client {
	cli := &client{}
	c := netc.New(ip, port, cli)
	c.Serve()
	return cli
}

func (this *client) OnConnect(conn net.Conn) {
	gcmd.Println(gcmd.Info, "连接到服务器:", conn.RemoteAddr())
}

func (this *client) OnDisConnect(conn net.Conn, reason string) {
	gcmd.Println(gcmd.Info, "与服务器断开连接:", conn.RemoteAddr(), "原因:", reason)
}

func (this *client) OnRecvMessage(conn net.Conn, msg []byte) {
	this.msgRecvHanlder(msg)

	// p := &api.Msg{}
	// p.Cmd.Command = 1
	// p.Cmd.Filepath = "file"
	// p.Cmd.Md5 = 75489259473
	// m := &api.Messages{}
	// m.Field1 = "1"
	// m.Field1 = "2"
	// m.Field1 = "3"
	// m.Field1 = "4"
	// p.Msg = append(p.Msg, m)

	// this.msgSendHanlder(p)
}

func (this *client) OnSendMessage(conn net.Conn) {

}

func (this *client) OnLocalCommand(conn net.Conn, cmd []byte, msg []byte) {
	gcmd.Println(gcmd.Ok, "本地命令:", string(cmd), string(msg))
}

func (this *client) Stop() {

}

func (this *client) msgRecvHanlder(msg []byte) {
	p := &api.Msg{}
	err := proto.Unmarshal(msg, p)
	if err != nil {
		gcmd.Println(gcmd.Err, err)
	}

	gcmd.Println(gcmd.Ok, "client recv:",p)
}

func (this *client) msgSendHanlder(msg *api.Msg) {
	b, err := proto.Marshal(msg)
	if err != nil {
		gcmd.Println(gcmd.Err, err)
	}

	gcmd.Println(gcmd.Ok, b)
}

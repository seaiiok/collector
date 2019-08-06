package netserver

import (
	"collector/api"
	"gcom/gcmd"
	"net"

	"github.com/golang/protobuf/proto"
	"snet/snet.v4"
)

type server struct {
}

func NewServer(ip string, port int) *server {
	svr := &server{}
	s := snet.NewServer(ip, port, svr)
	s.Serve()
	return svr
}

func (s *server) OnConnect(conn *net.TCPConn) {
	gcmd.Println(gcmd.Info, "客户端连接:", conn.RemoteAddr())
	p := &api.Msg{}
	c := &api.Commands{}
	c.Command = 1
	c.Filepath = "file"
	c.Md5 = 75489259473

	p.Cmd = c

	m := &api.Messages{}
	m.Field1 = "1"
	m.Field2 = "2"
	m.Field3 = "3"
	m.Field4 = "4"
	
	for i := 0; i < 10; i++ {
		p.Msg = append(p.Msg, m)
	}


	s.msgSendHanlder(conn, p)
}

func (s *server) OnDisConnect(conn *net.TCPConn, reason string) {
	gcmd.Println(gcmd.Info, "客户端断开连接:", conn.RemoteAddr(), "原因:", reason)
}

func (s *server) OnRecvMessage(conn *net.TCPConn, msg []byte) {
	s.msgRecvHanlder(conn, msg)
}

func (s *server) OnSendMessage(conn *net.TCPConn) {

}

func (s *server) OnLocalCommand(conn *net.TCPConn, cmd []byte, msg []byte) {
	gcmd.Println(gcmd.Ok, "本地命令:", string(cmd), string(msg))
}

func (s *server) Stop() {

}

func (s *server) msgRecvHanlder(conn *net.TCPConn, msg []byte) {
	p := &api.Msg{}
	err := proto.Unmarshal(msg, p)
	if err != nil {
		gcmd.Println(gcmd.Err, err)
	}

	gcmd.Println(gcmd.Ok, p)

}

func (s *server) msgSendHanlder(conn *net.TCPConn, msg *api.Msg) {
	b, err := proto.Marshal(msg)
	if err != nil {
		gcmd.Println(gcmd.Err, err)
	}

	gcmd.Println(gcmd.Ok, b)

	p := &snet.Package{}
	p.SetMsg(b)
	sendmsg := p.Pack()

	conn.Write(sendmsg)
}

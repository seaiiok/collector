package nets

import (
	"collector/api"
	"gcom/gcmd"
	"net"

	"snet/snet.v4"
)

type server struct {
}

func NewServer(ip, port string) *server {
	svr := &server{}
	s := snet.NewServer(ip, port, svr)
	s.Serve()
	return svr
}

func (s *server) OnConnect(conn *net.TCPConn) {
	gcmd.Println(gcmd.Info, "客户端连接:", conn.RemoteAddr())
}

func (s *server) OnDisConnect(conn *net.TCPConn, reason string) {
	gcmd.Println(gcmd.Info, "客户端断开连接:", conn.RemoteAddr(), "原因:", reason)
}

func (s *server) OnRecvMessage(conn *net.TCPConn, msg []byte) {
	s.msgRecvHanlder(msg)

	p := &api.Msg{}
	p.Cmd.Command = 1
	p.Cmd.Filepath = "file"
	p.Cmd.Md5 = 75489259473
	m := &api.Messages{}
	m.Field1 = "1"
	m.Field1 = "2"
	m.Field1 = "3"
	m.Field1 = "4"
	p.Msg = append(p.Msg, m)

	s.msgSendHanlder(p)
}

func (s *server) OnSendMessage(conn *net.TCPConn) {

}

func (s *server) OnLocalCommand(conn *net.TCPConn, cmd []byte, msg []byte) {
	gcmd.Println(gcmd.Ok, "本地命令:", string(cmd), string(msg))
}

func (s *server) Stop() {

}

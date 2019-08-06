package interfaces

import "net"

type INetClient interface {
	OnConnect(net.Conn)
	OnDisConnect(net.Conn, string)
	OnRecvMessage(net.Conn, []byte)
	OnSendMessage(net.Conn)
}

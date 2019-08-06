package netc

import (
	"context"
	"fmt"
	"gcom/gcmd"
	"net"
)

type client struct {
	ip     string
	port   int
	ctx    context.Context
	cancel context.CancelFunc
	// conn   net.Conn
	netClient INetClient
}

type INetClient interface {
	OnConnect(net.Conn)
	OnDisConnect(net.Conn, string)
	OnRecvMessage(net.Conn, []byte)
	OnSendMessage(net.Conn)
}

func New(ip string, port int, netClient INetClient) *client {
	ctx, cancel := context.WithCancel(context.Background())
	return &client{
		ip:        ip,
		port:      port,
		ctx:       ctx,
		cancel:    cancel,
		netClient: netClient,
	}
}

func (this *client) Serve() {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", this.ip, this.port))
	if err != nil {
		this.netClient.OnDisConnect(conn, err.Error())
		gcmd.Println(gcmd.Warn, "client dial err, exit!")
		return
	}

	this.netClient.OnConnect(conn)
	// c.conn = conn

	go this.newConnection(conn)

	for {
		select {
		case <-this.ctx.Done():
			defer conn.Close()
			return
		}
	}
}

func (this *client) Stop() {
	this.cancel()
}

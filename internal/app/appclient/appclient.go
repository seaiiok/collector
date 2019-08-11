package appclient

import (
	"collector/api"
	"collector/pkg/interfaces"
	"gcom/garchive"
	"net"
	"time"

	"github.com/golang/protobuf/proto"
	"snet/snet.v4/clients"
	"snet/snet.v4/packet"
)

type client struct {
	config    interfaces.IConfig
	producter interfaces.IProducter
	cache     interfaces.ICache
	log       interfaces.ILog
}
var netc interfaces.INetClient

func (this *client) OnConnect(conn net.Conn) {
	this.log.Info("连接到服务器:", conn.RemoteAddr().String())
}

func (this *client) OnDisConnect(conn net.Conn, reason string) {
	this.log.Info("断开与服务器连接:", conn.RemoteAddr().String(), ",原因:", reason)

	go func() {
		for {
			err := netc.Serve()
			if err != nil {
				time.Sleep(30 * time.Second)
				continue
			}
			break
		}
	}()
}

func (this *client) OnRecvMessage(conn net.Conn, msg []byte) {

	m := &api.Msg{}
	p := &packet.Package{}

	err := proto.Unmarshal(msg, m)
	if err != nil {
		this.log.Error("proto unmarshal err:", err)
		return
	}

	switch m.Cmd {
	case 0:

		this.producter.Set(m.File, []byte{1})

	case 1:
		for {
			select {
			case file := <-this.producter.OutQueue():
				m.Cmd = 2
				m.File = file

				pm, err := proto.Marshal(m)
				if err != nil {
					this.log.Error("proto marshal err:", err)
					this.producter.Set(m.File, []byte{1})
					return
				}

				pp, err := p.Pack(pm)
				if err != nil {
					this.log.Error("pack err:", err)
					this.producter.Set(m.File, []byte{1})
					return
				}

				conn.Write(pp)
			default:
				time.Sleep(1 * time.Second)
			}
		}

	case 3:
		b, md5, err := garchive.UnZip(m.File)
		if err != nil {
			this.log.Error("garchive unzip err:", err)
			this.producter.Set(m.File, []byte{1})
			return
		}

		m.Cmd = 4
		m.Md5 = md5
		m.Msg = b

		pm, err := proto.Marshal(m)
		if err != nil {
			this.log.Error("proto marshal err:", err)
			this.producter.Set(m.File, []byte{1})
			return
		}

		pp, err := p.Pack(pm)
		if err != nil {
			this.log.Error("pack err:", err)
			this.producter.Set(m.File, []byte{1})
			return
		}

		conn.Write(pp)

	case 5:
		this.producter.Set(m.File, []byte{1, 1})
	default:
	}

}

func (this *client) OnSendMessage(conn net.Conn) {

}


func NewClient(config interfaces.IConfig, producter interfaces.IProducter, cache interfaces.ICache, log interfaces.ILog) *client {

	c := &client{
		config:    config,
		producter: producter,
		cache:     cache,
		log:       log,
	}

	netc = snetclient.NewClient(config.Get("host"), config.Get("port"), c)

	go func() {
		for {
			err := netc.Serve()
			if err != nil {
				time.Sleep(30 * time.Second)
				continue
			}
			break
		}
	}()

	producter.Run()

	return c
}

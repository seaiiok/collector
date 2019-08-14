package appclient

import (
	"collector/api"
	"collector/pkg/interfaces"
	"fmt"
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

	//请求1-文件基础信息
	case 101:
		for {
			time.Sleep(5 * time.Second)
			select {
			case file := <-this.producter.OutQueue():

				list, _ := this.cache.GetMap()
				for k, v := range list {
					fmt.Println(k, "---", v)
				}
				fmt.Println()
				
				m.Cmd = 101
				m.File = file

				pm, err := proto.Marshal(m)
				if err != nil {
					this.log.Error("proto marshal err:", err)
					this.producter.Set(m.File, []byte{1})
					continue
				}

				pp, err := p.Pack(pm)
				if err != nil {
					this.log.Error("pack err:", err)
					this.producter.Set(m.File, []byte{1})
					continue
				}

				conn.Write(pp)
				return
			default:
				time.Sleep(1 * time.Second)
			}
		}

		//请求1-成功,请求2-发送文件完整信息
	case 102:

		b, md5, err := garchive.UnZip(m.File)
		if err != nil {
			this.log.Error("garchive unzip err:", err)
			this.producter.Set(m.File, []byte{1})
			return
		}

		m.Cmd = 102
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

		//请求1-失败,继续
	case 103:
		this.producter.Set(m.File, []byte{1})
		b := this.objectNext(m)
		conn.Write(b)

		//请求2-成功,继续
	case 104:
		this.producter.Set(m.File, []byte{1, 1})
		b := this.objectNext(m)
		conn.Write(b)

		//请求2-失败,继续
	case 105:
		this.producter.Set(m.File, []byte{1})
		b := this.objectNext(m)
		conn.Write(b)

		//请求1,2-拒绝处理,继续
	case 106:
		this.producter.Set(m.File, []byte{1, 1})
		b := this.objectNext(m)
		conn.Write(b)
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

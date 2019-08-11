package appserver

import (
	"collector/api"
	"collector/pkg/interfaces"
	"fmt"
	"gcom/garchive"
	"gcom/gtools/gdb"
	"net"
	"os"
	"strings"
	"time"

	"github.com/golang/protobuf/proto"
	"snet/snet.v4"
	"snet/snet.v4/packet"
)

type server struct {
	ip     string
	port   string
	config interfaces.IConfig
	log    interfaces.ILog
}

func NewServer(config interfaces.IConfig, log interfaces.ILog) *server {
	svr := &server{
		ip:     config.Get("host"),
		port:   config.Get("port"),
		config: config,
		log:    log,
	}
	this := snet.NewServer(svr.ip, svr.port, svr)
	this.Serve()
	return svr
}

func (this *server) OnConnect(conn *net.TCPConn) {
	m := &api.Msg{}
	p := &packet.Package{}

	remote := strings.Split(conn.RemoteAddr().String(), ":")
	remoteip := ""
	if len(remote) == 2 {
		remoteip = remote[0]
	} else {
		remoteip = "127.0.0.1"
	}

	updtime := time.Now().Format("2006-01-02 15:04:05")
	devices := fmt.Sprintf(sql_devices, remoteip, updtime, remoteip, remoteip, updtime)

	err := gdb.Exec(devices)
	if err != nil {
		this.log.Error(err)
	}
	this.log.Info("客户端连接:", remoteip)

	m.Cmd = 1
	m.Ip = remoteip

	b, _ := proto.Marshal(m)
	msg, _ := p.Pack(b)
	conn.Write(msg)
}

func (this *server) OnDisConnect(conn *net.TCPConn, reason string) {
	this.log.Info("客户端断开连接:", conn.RemoteAddr().String(), ",原因:", reason)
}

func (this *server) OnRecvMessage(conn *net.TCPConn, msg []byte) {
	this.recvMessageHanlder(conn, msg)
}

func (this *server) OnSendMessage(conn *net.TCPConn) {

}

func (this *server) Stop() {

}

func (this *server) recvMessageHanlder(conn *net.TCPConn, msg []byte) {
	pm := &api.Msg{}
	err := proto.Unmarshal(msg, pm)
	if err != nil {
		this.log.Error(err)
		return
	}

	updtime := time.Now().Format("2006-01-02 15:04:05")

	switch pm.Cmd {
	case 2:
		existsqlstring := fmt.Sprintf(sql_fileexists, pm.Ip, pm.File)
		exist, err := gdb.Exist(existsqlstring)
		if err != nil {
			this.log.Error("db file is exist err:", err)
		}

		if exist == true {
			b := this.objectFinish(pm)
			conn.Write(b)
		} else {
			b := this.objectNotFinish(pm)
			conn.Write(b)
		}

	case 4:

		os.MkdirAll("./read-pool", 0744)
		newzip := "./read-pool/" + pm.Ip + "-" + pm.Md5 + ".zip"

		// 判断文件是否存在
		if _, err := os.Stat(newzip); os.IsExist(err) {
			b := this.objectFinish(pm)
			conn.Write(b)
			return
		}

		err := garchive.Zip(newzip, pm.Md5, pm.Msg)
		if err != nil {
			b := this.objectFailed(pm)
			conn.Write(b)
			defer func() {
				os.Remove(newzip)
			}()
			return
		}

		insertfiles := fmt.Sprintf(sql_insertfile, pm.Ip, pm.File, pm.Md5, updtime, "0")
		err = gdb.Exec(insertfiles)
		if err != nil {
			this.log.Error("db insert err:", err)
			return
		}
		pm.Cmd = 5
		pm.Msg = []byte{}
		b := this.objectFinish(pm)
		conn.Write(b)
	default:

	}
}

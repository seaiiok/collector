package appserver

import (
	"collector/api"
	"collector/pkg/interfaces"
	"context"
	"fmt"
	"gcom/gtools/gdb"
	"io/ioutil"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"
	"snet/snet.v4"
)

var mu sync.Mutex

type server struct {
	ip     string
	port   string
	config interfaces.IConfig
	log    interfaces.ILog
	ctx    context.Context
}

func NewServer(ctx context.Context, config interfaces.IConfig, log interfaces.ILog) *server {
	svr := &server{
		ip:     config.Get("host"),
		port:   config.Get("port"),
		config: config,
		log:    log,
		ctx:    ctx,
	}
	this := snet.NewServer(svr.ip, svr.port, svr)
	this.Serve()

	//关闭net serve
	go func() {
		for {
			select {
			case <-ctx.Done():
				this.Stop()
				return
		
			}
		}
	}()

	return svr
}

func (this *server) OnConnect(conn *net.TCPConn) {
	pm := &api.Msg{}

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
		this.log.Error(err, ";this client service stop!")
		return
	}
	this.log.Info("客户端连接:", remoteip)

	pm.Ip = remoteip
	b := this.objectNext(pm, 101)
	conn.Write(b)
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
	case 100:
		b := this.objectNext(pm, 101)
		conn.Write(b)

	case 101:
		existsqlstring := fmt.Sprintf(sql_fileexists, pm.Ip, pm.File)
		exist, err := gdb.Exist(existsqlstring)
		if err != nil {
			b := this.objectNext(pm, 103)
			conn.Write(b)
			this.log.Error("db file exists err:", err, ";this file:", pm.Ip,pm.File)
			return
		}

		if exist == true {
			b := this.objectNext(pm, 106)
			conn.Write(b)
			this.log.Info("this file is exists in db:", pm.Ip,pm.File)
		} else {
			b := this.objectNext(pm, 102)
			conn.Write(b)
		}

	case 102:
		newfile := readPool + pm.Ip + "-" + pm.Md5 + ".txt"
		fFinish := "0"
		// 判断文件是否存在
		if _, err := os.Stat(newfile); err == nil {
			// os.IsNotExist(err)
			b := this.objectNext(pm, 106)
			conn.Write(b)
			this.log.Info("this file is exists in read-pool:", pm.Ip,pm.File)
			return
		}

		err := ioutil.WriteFile(newfile, pm.Msg, 0644)
		if err != nil {
			defer func() {
				os.Remove(newfile)
			}()

			b := this.objectNext(pm, 105)
			conn.Write(b)
			this.log.Error("write file to read-pool err:", err, ";this file:", pm.Ip,pm.File)
			return
		}

		insertfiles := fmt.Sprintf(sql_insertfile, pm.Ip, pm.File, pm.Md5, updtime, fFinish)
		err = gdb.Exec(insertfiles)
		if err != nil {
			defer func() {
				os.Remove(newfile)
			}()

			b := this.objectNext(pm, 105)
			conn.Write(b)
			this.log.Error("db insert file err:", err, ";this file:", pm.Ip,pm.File)
			return
		}

		b := this.objectNext(pm, 104)
		_, err = conn.Write(b)
		if err != nil {
			return
		}
		this.log.Info("this file go read-pool:", pm.Ip,pm.File)
	default:
		conn.Write([]byte("unknown request,plz format msg and send again!"))
	}
}

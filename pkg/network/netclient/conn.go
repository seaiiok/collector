package netc

import (
	"io"
	"net"

	"snet/snet.v4"
)

func (this *client) newConnection(conn net.Conn) {
	go this.netClient.OnSendMessage(conn)

	for {
		select {
		case <-this.ctx.Done():
			defer conn.Close()
			return
		default:
			this.remoteConnHandle(conn)
		}
	}
}

func (this *client) remoteConnHandle(conn net.Conn) {
	p := &snet.Package{}
	tempBuff := make([]byte, 0)

	for {

		select {
		case <-this.ctx.Done():
			defer conn.Close()
			return
		default:
		}

		msg := make([]byte, 0)
		oneByte := make([]byte, 1)
		n, err := conn.Read(oneByte)
		if err != nil {
			if err == io.EOF {
				//TODO client close
				this.netClient.OnDisConnect(conn, "client connection close")
				this.Stop()
				return
			}
			this.netClient.OnDisConnect(conn, err.Error())
			this.Stop()
			return
		}

		if n != 1 {
			continue
		}

		tempBuff = append(tempBuff, oneByte...)
		tempBuffLength := len(tempBuff)
		if tempBuffLength < 8 {
			continue
		}

		tempBuff = tempBuff[tempBuffLength-8:]
		tempMsgLength := p.UnPackMsgLength(tempBuff[:4])
		tempMsgLengthCRC1 := p.CheckCRC32(tempMsgLength)
		tempMsgLengthCRC2 := p.UnPackMsgLength(tempBuff[4:8])
		if tempMsgLengthCRC1 != tempMsgLengthCRC2 {
			continue
		}

		msgBytes := make([]byte, tempMsgLength)
		n, err = io.ReadFull(conn, msgBytes)
		if err != nil {
			if err == io.EOF {
				//TODO client close
				this.netClient.OnDisConnect(conn, "client connection close")
				this.Stop()
				return
			}
			this.netClient.OnDisConnect(conn, err.Error())
			this.Stop()
			return
		}

		if n != int(tempMsgLength) {
			continue
		}

		msg = append(msg, msgBytes...)
		this.netClient.OnRecvMessage(conn, msg)
		return

	}
}

package appserver

import (
	"collector/api"

	"github.com/golang/protobuf/proto"
	"snet/snet.v4/packet"
)

func (this *server) objectFinish(pm *api.Msg) []byte {
	pack := &packet.Package{}
	pm.Cmd = 5
	pm.Msg = []byte{}
	b, err := proto.Marshal(pm)
	if err != nil {
		return []byte{}
	}
	msg, err := pack.Pack(b)
	if err != nil {
		return []byte{}
	}
	return msg
}

func (this *server) objectNotFinish(pm *api.Msg) []byte {
	pack := &packet.Package{}
	pm.Cmd = 3
	pm.Msg = []byte{}
	b, err := proto.Marshal(pm)
	if err != nil {
		return []byte{}
	}
	msg, err := pack.Pack(b)
	if err != nil {
		return []byte{}
	}
	return msg
}

func (this *server) objectFailed(pm *api.Msg) []byte {
	pack := &packet.Package{}
	pm.Cmd = 0
	pm.Msg = []byte{}
	b, err := proto.Marshal(pm)
	if err != nil {
		return []byte{}
	}
	msg, err := pack.Pack(b)
	if err != nil {
		return []byte{}
	}
	return msg
}

package appclient

import (
	"collector/api"

	"github.com/golang/protobuf/proto"
	"snet/snet.v4/packet"
)

func (this *client) objectNext(pm *api.Msg) []byte {
	pack := &packet.Package{}
	pm.Cmd = 100
	pm.Msg = []byte{}
	pm.File = ""
	pm.Md5 = ""

	b, err := proto.Marshal(pm)
	if err != nil {
		this.log.Error(err, "; service stop!")
		return []byte{}
	}
	msg, err := pack.Pack(b)
	if err != nil {
		this.log.Error(err, "; service stop!")
		return []byte{}
	}
	return msg
}

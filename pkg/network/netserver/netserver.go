package nets

import (
	"collector/api"
	"gcom/gcmd"

	"github.com/golang/protobuf/proto"
)

func (s *server) msgRecvHanlder(msg []byte) {
	p := &api.Msg{}
	err := proto.Unmarshal(msg, p)
	if err != nil {
		gcmd.Println(gcmd.Err, err)
	}

	gcmd.Println(gcmd.Ok, p)
}

func (s *server) msgSendHanlder(msg *api.Msg) {
	b, err := proto.Marshal(msg)
	if err != nil {
		gcmd.Println(gcmd.Err, err)
	}

	gcmd.Println(gcmd.Ok, b)
}

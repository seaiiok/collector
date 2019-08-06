package api

import (
	"testing"

	proto "github.com/golang/protobuf/proto"
)

func TestApi(t *testing.T) {
	// u1 := &Msg{
	// 	Cmd: &Command{
	// 		Cmd: 99,
	// 	},
	// 	Msg: {
	// 		{Id: 101, Name: "Tom1", Desc: "三好学生1"},
	// 		{Id: 102, Name: "Tom2", Desc: "三好学生2"},
	// 		{Id: 103, Name: "Tom3", Desc: "三好学生3"},
	// 	},
	// }
	m1 := NewMsg()

	m1.Cmd.Command = 99
	m1.Cmd.Filepath = "fhjasjkj"

	for i := 0; i < 10; i++ {
		m0 := &Messages{Field1: "101", Field2: "Tom", Field3: "三好学生1"}
		m1.Msg = append(m1.Msg, m0)
	}

	data, _ := proto.Marshal(m1)
	t.Log("序列化长度", len(data))
	m2 := NewMsg()
	proto.Unmarshal(data, m2)

	t.Log(m2)

}

func NewMsg() *Msg {
	m := &Msg{}
	m.Cmd = &Commands{}
	m.Msg = make([]*Messages, 0)
	return m
}

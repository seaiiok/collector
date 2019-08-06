package server

// //获取数据
// func MsgProducer(k string, v int) {
// 	gcmd.Println(gcmd.Ok, "消息文件读取中:", k)
// 	b, err := gfiles.ReadZip(k)
// 	if err != nil {
// 		gcmd.Println(gcmd.Err, err)
// 	}

// 	msgs := bytes.Split(b, []byte("\n"))

// 	for _, msg := range msgs {
// 		if len(msg) < 69 || len(msg) > 70 || v != 127 {
// 			continue
// 		}

// 		p.MsgQueue <- msg
// 		localcache.PutCache(k, v+gfiles.SCAN_READED)
// 	}

// }

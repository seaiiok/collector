package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"

	"github.com/mholt/archiver"
	"github.com/seaiiok/gcom/gfiles" 
)

func main() {
	// go GoServer()
	// go GoClient()

	// select{}
	file := "D:/FTP/online/2019/1.zip"
	file2 := "D:/FTP/online/2019/1.log"
	b, err := gfiles.DeCompressZip(file)
	if err != nil {
		fmt.Println(err)
	}
	// ""

	m := md5.New()
	m.Write(b)

	// sum:=hex.EncodeToString(m.Sum(nil))

	// fmt.Println(sum)

	md5str2 := fmt.Sprintf("%x", m.Sum(nil))

	fmt.Println(md5str2)

	res, err := GetFileMd5(file2)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res)

	err = archiver.Archive([]string{"testdata", file2}, "test.zip")
	if err != nil {
		fmt.Println(err)
	}
	// err = archiver.Unarchive("test.tar.gz", "test")
}

// func GoServer() {
// 	netserver.NewServer("127.0.0.1", 496)
// }

// func GoClient() {
// 	netclient.NewClient("127.0.0.1", 496)
// }

func GetFileMd5(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("os Open error")
		return "", err
	}
	md5 := md5.New()
	_, err = io.Copy(md5, file)
	if err != nil {
		fmt.Println("io copy error")
		return "", err
	}
	md5Str := hex.EncodeToString(md5.Sum(nil))
	return md5Str, nil
}

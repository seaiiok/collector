package fileparse

import (
	"bufio"
	"bytes"
	"gcom/gtools/gdb"
	"io"
	"os"
)

const (
	icount   = 5000
	finished = "file finished"
	MAXFILES = 100
	MAXINFOS = 5000
)

type files struct {
	fileIp   string
	fileMd5  string
	progress string
}

type infos struct {
	id string
}

var filesChan = make(chan files, MAXFILES)
var infosChan = make(chan infos, MAXINFOS)

func Run() {

}

func run() {
	fs := make([]files, 0)
	for {
		select {
		default:
			if len(filesChan) == 0 && len(infosChan) == 0 {
				fs, _ = getFilesName()
				for _, v := range fs {
					if len(filesChan) >= MAXFILES {
						break
					}
					filesChan <- v
				}
			}
		}
	}
}

//获取数据库未完成的文件信息
func getFilesName() (fs []files, err error) {
	fs = make([]files, 0)
	res, err := gdb.Querys(sql_notfinishfile)
	if err != nil {
		return fs, err
	}

	for _, v := range res {
		f := files{}
		if len(v) != 3 {
			continue
		}

		f.fileIp = v[0]
		f.fileMd5 = v[1]
		f.progress = v[2]
		fs = append(fs, f)
	}
	return fs, nil
}

func ParseFile(file string) ([][]byte, error) {
	res := make([][]byte, 0)

	f, err := os.OpenFile(file, os.O_RDONLY, 0644)
	if err != nil {
		return res, err
	}

	r := bufio.NewReader(f)

	scan := bufio.NewScanner(r)

	scan.Split(bufio.ScanLines)

	for scan.Scan() {
		sb := bytes.TrimPrefix(scan.Bytes(), []byte{239, 187, 191})
		res = append(res, sb)
	}
	err = scan.Err()
	return res, err
}

func ParseRows(file string, seek int64) ([][]string, error) {
	rows := make([][]string, 0)
	bs, err := ParseFile(file)
	if len(bs) <= int(seek) {
		return nil, io.EOF
	}
	if err != nil {
		return nil, err
	}

	bs = bs[seek:]

	for _, v := range bs {
		if len(v) > 0 {
			if v[0] == 32 {
				continue
			}
		}

		if len(v) < 69 {
			continue
		}
		r := make([]string, 0)
		r = append(r, string(v[:12]), string(v[13:23]), string(v[64:66]), string(v[67:69]))
		rows = append(rows, r)
	}
	return rows, err
}

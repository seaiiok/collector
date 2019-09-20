package fileparse

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"gcom/gtools/gdb"
	"os"
	"strconv"
	"time"
)

const (
	icount = 5000
)

type files struct {
	fileIp   string
	fileMd5  string
	progress string
}

func Run(ctx context.Context, root string) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				err := run(root)
				if err != nil {
					time.Sleep(10 * time.Second)
					fmt.Println("run err:", err)
					continue
				}
			}
		}
	}()
}

func getNotFinishFiles() (files, err) {
	f := files{}

	res, err := gdb.Querys(sql_notfinishfile)
	if err != nil {
		return files{}, err
	}

	if len(res) == 0 {
		time.Sleep(10 * time.Second)
		return files{}, err
	}

	for _, v := range res {
		if len(v) != 3 {
			continue
		}
		f.fileIp = v[0]
		f.fileMd5 = v[1]
		f.progress = v[2]
	}
	return f, err
}

func run(root string) error {

	res, err := gdb.Querys(sql_notfinishfile)
	if err != nil {
		return err
	}

	if len(res) != 1 {
		time.Sleep(10 * time.Second)
		return err
	}

	if len(res[0]) != 3 {
		return err
	}

	fileip := res[0][0]
	filemd5 := res[0][1]
	progress := res[0][2]
	seek, err := strconv.ParseInt(progress, 10, 64)
	if err != nil {
		seek = 0
	}
	updtime := time.Now().Format("2006-01-02 15:04:05")
	file := root + fileip + "-" + filemd5 + ".txt"
	rows, err := ParseRows(file, seek)

	if err != nil {
		return err
	}
	var upseekInt = int(seek)

	irows := make([][]string, 0)
	for _, v := range rows {
		nv := make([]string, 0)
		nv = append(nv, filemd5)
		nv = append(nv, v...)
		nv = append(nv, updtime)
		irows = append(irows, nv)
	}

	for {
		lenIrows := len(irows)
		if lenIrows == 0 {
			break
		}

		if lenIrows < 5000 {
			err = gdb.Insertbulk(sql_insertinfo, irows)
			if err != nil {
				break
			}

			irows = make([][]string, 0)
			upseekInt = upseekInt + lenIrows
			break
		}

		if lenIrows >= 5000 {
			err = gdb.Insertbulk(sql_insertinfo, irows[:5000])
			if err != nil {
				break
			}
			irows = irows[5000:]
			upseekInt = upseekInt + 5000
			continue
		}
	}

	if err != nil {
		return err
	}

	upseek := strconv.Itoa(upseekInt)
	upfinish := "1"
	upd := []string{upseek, upfinish, updtime, filemd5, fileip}

	upds := make([][]string, 0)

	upds = append(upds, upd)
	err = gdb.Updates(sql_updateprogress, upds)

	return err
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
	if len(bs) <= int(seek) || err != nil {
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
		//fmt.Println(v)
		//fmt.Println(string(v[:12]), string(v[13:23]), string(v[64:66]), string(v[67:69]))
		rows = append(rows, r)
	}
	return rows, err
}

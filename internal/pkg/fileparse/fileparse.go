package fileparse

import (
	"bufio"
	"bytes"
	"context"
	"gcom/gtools/gdb"
	"os"
	"strconv"
	"time"
)

const (
	icount = 5000
)

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
					continue
				}
			}
		}
	}()
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
	rows, n, err := ParseRows(file, seek)

	if err != nil {
		return err
	}

	var ct int64
	for {

		irows := make([][]string, 0)
		var i int64
		for _, v := range rows[ct:] {
			ct++
			i++
			nv := make([]string, 0)
			nv = append(nv, filemd5)
			nv = append(nv, v...)
			nv = append(nv, updtime)
			irows = append(irows, nv)

			//限制插入数据量
			if i >= icount {
				break
			}
		}

		err = gdb.Insertbulk(sql_insertinfo, irows)

		if err != nil {
			return err
		}

		upseek := strconv.FormatInt(ct, 10)
		upfinish := "0"
		if len(irows) == 0 || ct == n {
			upfinish = "1"
		}

		upd := []string{upseek, upfinish, updtime, filemd5, fileip}

		upds := make([][]string, 0)

		upds = append(upds, upd)
		err = gdb.Updates(sql_updateprogress, upds)

		if upfinish == "1" {
			return nil
		}
	}

}

func ParseFile(file string, seek int64) ([][]byte, int64, error) {
	var lines int64
	res := make([][]byte, 0)

	f, err := os.OpenFile(file, os.O_RDONLY, 0644)
	if err != nil {
		return res, 0, err
	}

	r := bufio.NewReader(f)

	scan := bufio.NewScanner(r)

	scan.Split(bufio.ScanLines)

	for scan.Scan() {
		lines++
		if lines > seek {
			sb := bytes.TrimPrefix(scan.Bytes(), []byte{239, 187, 191})
			res = append(res, sb)
		}
	}
	err = scan.Err()
	return res, lines, err
}

func ParseRows(file string, seek int64) ([][]string, int64, error) {
	rows := make([][]string, 0)
	bs, n, err := ParseFile(file, seek)
	if len(bs) == 0 || err != nil {
		return nil, n, err
	}

	for _, v := range bs {
		if len(v) < 69 {
			continue
		}
		r := make([]string, 0)
		r = append(r, string(v[:12]), string(v[13:23]), string(v[64:66]), string(v[67:69]))

		rows = append(rows, r)
	}
	return rows, n, err
}

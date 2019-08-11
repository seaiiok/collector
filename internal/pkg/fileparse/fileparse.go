package fileparse

import (
	"archive/zip"
	"bufio"
	"bytes"
	"fmt"
	"gcom/gtools/gdb"
	"io"
	"os"
	"strconv"
	"time"
)

func Run() {
	notFinishFile := fmt.Sprintf(`SELECT top 1 Files,MD5,Progress FROM [iFixsvr_JFOffline_info] WHERE Finish <> '1'`)
	res, err := gdb.Querys(notFinishFile)
	if err != nil {
		fmt.Println("err1:", err)
	}

	if len(res[0]) != 3 {
		fmt.Println("结果:", len(res))
	}
	file := res[0][0]
	progress := res[0][2]
	seek, err := strconv.ParseInt(progress, 10, 64)
	if err != nil {
		fmt.Println("err2:", err)
		seek = 0
	}
	fmt.Println(file, seek)
	rows, n, err := ParseRows(file, seek)

	fmt.Println("----len rows", len(rows))
	updateProgress := `UPDATE iFixsvr_JFOffline_info SET Progress = ?, Finish = ?, FinishDTime = ? WHERE (Files = ?)`

	upseek := strconv.FormatInt(n, 10)
	upfinish := "0"
	if err == io.EOF {
		upfinish = "1"
	}

	updtime := time.Now().Format("2006-01-02 15:04:05")
	upfile := file
	fmt.Println(upseek, upfinish, updtime, upfile)
	upd := []string{upseek, upfinish, updtime, upfile}

	upds := make([][]string, 0)

	upds = append(upds, upd)
	err = gdb.Updates(updateProgress, upds)
	fmt.Println(err)
}

func ParseFile(file string, seek int64) ([][]byte, int64, error) {
	var lines int64
	res := make([][]byte, 0)
	f, err := os.OpenFile(file, os.O_RDONLY, 0744)
	if err != nil {
		return nil, 0, nil
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

func ParseFile1(path string, seek int64) ([][]byte, int64, error) {
	var lines int64
	rc, err := zip.OpenReader(path)
	defer rc.Close()
	bs := make([][]byte, 0)

	if err != nil {
		return bs, 0, err
	}

	for _, file := range rc.File {
		f, err := file.Open()
		if err != nil {
			return bs, 0, err
		}

		r := bufio.NewReader(f)

		scan := bufio.NewScanner(r)

		scan.Split(bufio.ScanLines)

		for scan.Scan() {
			lines++
			if lines > seek {
				sb := bytes.TrimPrefix(scan.Bytes(), []byte{239, 187, 191})
				bs = append(bs, sb)
			}
		}
		err = scan.Err()
		return bs, lines, err
	}
	return bs, 0, err
}

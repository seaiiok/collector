package main

import (
	"collector/internal/app/appclient"
	"collector/pkg/cache"
	"collector/pkg/config"
	"collector/pkg/global"
	"collector/pkg/logger"
	"errors"
	"fmt"
	"gcom/gwin"
	"io/ioutil"
	"os"
	"strconv"
)

const (
	logroot        = "./log"
	loglevel       = logger.LOGALL
	cachefile      = "./cache"
	cachebucket    = "db"
	configfile     = "./configs/collector.json"
	apppid         = "mescollectclient.pid"
	messagecaption = "程序已经启动!"
)

func init() {
	iManPid := fmt.Sprint(os.Getpid())
	tmpDir := os.TempDir()

	if err := procExsit(tmpDir); err == nil {
		pidFile, _ := os.Create(tmpDir + "\\" + apppid)
		defer pidFile.Close()

		pidFile.WriteString(iManPid)
	} else {
		gwin.MessageBox(global.MESSAGEBOXCAPTION, messagecaption, gwin.MB_ICONWARNING|gwin.MB_OK)
		os.Exit(1)
	}
}

// 判断进程是否启动
func procExsit(tmpDir string) (err error) {
	iManPidFile, err := os.Open(tmpDir + "\\" + apppid)
	defer iManPidFile.Close()

	if err == nil {
		filePid, err := ioutil.ReadAll(iManPidFile)
		if err == nil {
			pidStr := fmt.Sprintf("%s", filePid)
			pid, _ := strconv.Atoi(pidStr)
			_, err := os.FindProcess(pid)
			if err == nil {
				return errors.New("app online")
			}
		}
	}

	return nil
}

func main() {

	app := appclient.New(logger.New(logroot, loglevel), cache.New(cachefile, cachebucket), config.New(configfile))
	app.Serve()
	app.Run()
}

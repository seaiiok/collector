// user "\r\n"
// +build only windows

package logger

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	logname  = "log "
	enter    = "\r\n\r\n"
	LOGALL   = 0
	LOGINFO  = 1
	LOGWARN  = 2
	LOGERROR = 3
)

var (
	lastmouth = "200601"
	lastday   = "2006-01-02"
	once      = sync.Once{}
	file      = &os.File{}
	sw        = &sync.RWMutex{}
)

type logger struct {
	path  string
	level int
}

func New(path string, level int) *logger {
	return &logger{
		path:  path,
		level: level,
	}
}

func (this *logger) init() {
	sw.Lock()
	defer sw.Unlock()

	pmouth := time.Now().Format("200601")

	if !strings.EqualFold(pmouth, lastmouth) {
		file.Close()

		extime := time.Now().Format("200601")
		logfile := this.path + "/" + logname + extime + ".log"
		if _, err := os.Stat(this.path); !os.IsExist(err) {
			os.MkdirAll(this.path, 0644)
		}

		file, _ = os.OpenFile(logfile, os.O_CREATE|os.O_APPEND, 0644)
		lastmouth = pmouth
	}

}

func (this *logger) Info(v ...interface{}) {
	if this.level >= LOGALL && this.level <= LOGINFO {
		this.init()
		logTextFormat("INFO", fmt.Sprintf("%s", v...))
	}
}

func (this *logger) Warn(v ...interface{}) {
	if this.level >= LOGALL && this.level <= LOGWARN {
		this.init()
		logTextFormat("WARN", fmt.Sprintf("%s", v...))
	}
}

func (this *logger) Error(v ...interface{}) {
	if this.level >= LOGALL && this.level <= LOGERROR {
		this.init()
		logTextFormat("ERRS", fmt.Sprintf("%s", v...))
	}
}

func logTextFormat(levelinfo string, v string) {
	sw.Lock()
	defer sw.Unlock()

	pday := time.Now().Format("2006-01-02")
	if !strings.EqualFold(pday, lastday) {
		file.WriteString("=========================== [" + pday + "] ===========================" + enter)
		lastday = pday
	}
	time := time.Now().Format("15:04:05.000000")
	file.WriteString("[" + levelinfo + "] " + time + " " + v + enter)
}

func (this *logger) Close() {
	if file != nil {
		file.Close()
	}
}

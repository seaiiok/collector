package log

import (
	"collector/pkg/interfaces"
	"gcom/glog"
)

type log struct {
	logroot string
	level   int
}

func New(logroot string, level int) interfaces.ILog {
	this := &log{
		logroot: logroot,
		level:   level,
	}

	this.init()

	return this
}

func (this *log) init() {
	glog.Init(this.logroot, this.level)
}

func (this *log) Info(v ...interface{}) {
	glog.Info(v)
}

func (this *log) Warn(v ...interface{}) {
	glog.Warn(v)
}

func (this *log) Error(v ...interface{}) {
	glog.Error(v)
}

func (this *log) Close() {
	glog.Close()
}

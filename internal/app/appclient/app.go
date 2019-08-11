package appclient

import (
	"collector/internal/pkg/filequeue"
	"collector/pkg/global"
	"collector/pkg/interfaces"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type app struct {
	config    interfaces.IConfig
	producter interfaces.IProducter
	cache     interfaces.ICache
	log       interfaces.ILog
}

func New(config interfaces.IConfig, cache interfaces.ICache, log interfaces.ILog) *app {
	dirs := strings.Split(config.Get("collectroot"), global.SEP)
	if len(dirs) == 0 {
		log.Error("config file [collectroot] err! app exit!")
		panic("check log")
	}

	size, err := strconv.ParseInt(config.Get("filessize"), 10, 64)
	if err != nil {
		log.Error("config file [filessize] err! app exit!")
		panic("check log")
	}

	ext := strings.Split(config.Get("filesext"), global.SEP)
	if len(dirs) == 0 {
		log.Error("config file [filesext] err! app exit!")
		panic("check log")
	}

	cStr := strings.Split(config.Get("fileskeep"), global.SEP)
	if len(dirs) == 0 {
		log.Error("config file [fileskeep] err! app exit!")
		panic("check log")
	}

	eStr := strings.Split(config.Get("fileslose"), global.SEP)
	if len(dirs) == 0 {
		log.Error("config file [fileslose] err! app exit!")
		panic("check log")
	}

	files := filequeue.New(dirs, log, cache, &filter{size, ext, cStr, eStr})

	return &app{
		config:    config,
		producter: files,
		cache:     cache,
		log:       log,
	}
}

func (this *app) Run() {
	NewClient(this.config, this.producter, this.cache, this.log)
}

type filter struct {
	Size  int64
	Exts  []string
	CStrs []string
	EStrs []string
}

func (this *filter) Filter(list map[string]os.FileInfo) map[string]os.FileInfo {
	templist := make(map[string]os.FileInfo)
	f := false

	for k, v := range list {
		f = false
		if v.Size() < this.Size {
			continue
		}

		for _, ext := range this.Exts {
			if strings.EqualFold(filepath.Ext(v.Name()), ext) {
				f = true
				break
			}
		}

		if !f {
			continue
		}

		f = false
		for _, cstr := range this.CStrs {
			if strings.Contains(v.Name(), cstr) {
				f = true
				break
			}
		}

		if !f {
			continue
		}

		f = false
		for _, estr := range this.EStrs {
			if strings.Contains(v.Name(), estr) {
				f = true
				break
			}
		}

		if f {
			continue
		}

		templist[k] = v

	}

	return templist
}

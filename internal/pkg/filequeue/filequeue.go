package filequeue

import (
	"collector/pkg/global"
	"collector/pkg/interfaces"
	"context"
	"gcom/gfiles"
	"os"
	"time"
)

type IFilter interface {
	Filter(map[string]os.FileInfo) map[string]os.FileInfo
}

type filesQueue struct {
	dirs       []string
	filesqueue chan string
	queue      chan string
	//
	ctx    context.Context
	cancel context.CancelFunc

	//日志
	log interfaces.ILog
	//缓存
	cache interfaces.ICache
	//过滤接口
	filters IFilter
}

func New(dirs []string, log interfaces.ILog, cache interfaces.ICache, filters IFilter) interfaces.IProducter {

	fq := &filesQueue{
		dirs:       dirs,
		filesqueue: make(chan string, global.MAXQUEUEFILES),
		queue:      make(chan string, 1),
		//日志
		log: log,
		//缓存
		cache: cache,
		//过滤接口
		filters: filters,
	}
	return fq
}

func (this *filesQueue) Run() {
	ctx, cancel := context.WithCancel(context.Background())
	this.ctx = ctx
	this.cancel = cancel

	go func() {
		this.log.Info("文件队列开始工作...")
		list := make(map[string]os.FileInfo)
		for {
			select {
			case <-this.ctx.Done():
				return
			case <-time.After(10 * time.Second):
				list = make(map[string]os.FileInfo)

				//扫描出所有文件
				m, err := scanFiles(this.dirs...)
				if err != nil {
					continue
				}

				//过滤已处理过的文件

				for k := range m {
					m1, _ := this.cache.GetMap()
					if ok := m1[k]; len(ok) == 0 {
						f, err := os.Stat(k)
						if os.IsNotExist(err) {
							continue
						}
						list[k] = f
					}
				}

				// 过滤文件

				list = this.filters.Filter(list)

				//持久化
				for k := range list {
					this.Set(k, []byte{1})
					this.log.Info("加入文件队列:", k)
				}

				this.outQueue()
			}
		}
	}()
}

func (this *filesQueue) Stop() {
	this.log.Info("文件队列停止工作...")
	this.cancel()
}

//缓存文件
func (this *filesQueue) Set(k string, v []byte) {
	value := this.cache.Get(k)
	if len(value) > 3 {
		return
	}
	value = append(value, v...)
	this.cache.Set(k, value)
}

func (this *filesQueue) Get(k string) []byte {
	return this.cache.Get(k)
}

func (this *filesQueue) OutQueue() (queue chan string) {
	this.queue <- <-this.filesqueue
	return this.queue
}

func (this *filesQueue) outQueue() {
	if len(this.filesqueue) > 0 {
		return
	}

	m, err := this.cache.GetMap()
	if err != nil || len(m) == 0 {
		return
	}

	for k, v := range m {
		if len(this.filesqueue) >= global.MAXQUEUEFILES {
			break
		}

		if len(v) < 3 {
			this.filesqueue <- k
		}
	}
}

//scanFiles 遍历文件
func scanFiles(dirs ...string) (list map[string]struct{}, err error) {
	list = make(map[string]struct{})
	if len(dirs) == 0 {
		return
	}
	for _, dir := range dirs {
		fs, err := gfiles.ScanFiles(dir)
		if err != nil {
			return list, err
		}

		for k, v := range fs {
			list[k] = v
		}
	}
	return list, nil
}

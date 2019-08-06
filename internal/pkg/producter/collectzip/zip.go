package collectzip

import (
	"collector/internal/interfaces"
	"collector/pkg/global"
	"context"
	"sync"
	"time"
)

type zipBag struct {
	Files      map[string]struct{}
	QueueFiles chan string
	QueueMsg   chan []byte

	ctx    context.Context
	cancel context.CancelFunc

	mx *sync.RWMutex

	log       interfaces.ILog
	cache     interfaces.ICache
	config    interfaces.IConfig
	netclient interfaces.INetClient
}

func New(log interfaces.ILog, cache interfaces.ICache, config interfaces.IConfig, netclient interfaces.INetClient) *zipBag {
	ctx, cancel := context.WithCancel(context.Background())
	zip := &zipBag{
		Files:      make(map[string]struct{}, 0),
		QueueFiles: make(chan string, global.MAXQUEUEFILES),
		QueueMsg:   make(chan []byte, global.MAXQUEUEMSG),
		ctx:        ctx,
		cancel:     cancel,
		mx:         &sync.RWMutex{},
		log:        log,
		cache:      cache,
		config:     config,
		netclient:  netclient,
	}

	return zip
}

func (z *zipBag) Run() {
	go func() {
		for {
			select {
			case <-z.ctx.Done():
				return
			case <-time.After(60 * time.Second):

				if len(z.QueueFiles) >= 1 {
					continue
				}

				z.scanLocalFiles()

				z.scanFilesExt()

				z.scanFilesKeep()

				z.scanFilesLose()

				z.scanFilesSize()

				z.scanFilesExpired()

				z.expectFileCache()

				// 维持文件队列
				z.queueFiles()

			}
		}
	}()

}

func (z *zipBag) Stop() {
	z.cancel()
}

func (z *zipBag) Output() chan string {
	return z.QueueFiles
}

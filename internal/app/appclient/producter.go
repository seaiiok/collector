package appclient

import (
	"fmt"
	"time"
)

func (this *app) serve() {
	go func() {
		//文件采集器
		this.producter.Run()

		for {
			select {
			case <-this.ctx.Done():
				return
			case file := <-this.producter.Output():
				this.producter.DoneAFile(file)
				fmt.Println("QueueFiles:", file)
				fmt.Println()
				time.Sleep(1 * time.Second)
			}
		}
	}()

}

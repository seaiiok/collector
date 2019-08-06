package collectzip

import (
	"collector/pkg/global"
)

// 文件队列
func (z *zipBag) queueFiles() {

	if len(z.QueueFiles) >= 1 {
		return
	}

	files := z.cache.GetMap()
	for k, v := range files {
		if len(z.QueueFiles) >= global.MAXQUEUEFILES {
			return
		}
		if v.(string) == global.FILENOTREAD {
			z.QueueFiles <- k
		}
	}

}

// 文件完成
func (z *zipBag) DoneAFile(file string) {
	z.cache.Set(file, global.FILEREADED)
}

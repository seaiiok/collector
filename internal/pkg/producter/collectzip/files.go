package collectzip

import (
	"collector/pkg/global"
	"gcom/gfiles"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	ZipCollectRoot         = "collectroot"
	ZipFilesExt            = "filesext"
	ZipFilesKeep           = "fileskeep"
	ZipFilesLose           = "fileslose"
	ZipFilesSize           = "filessize"
	ZipFilesExpired        = "filesexpired"
	ZipFilesDefaultSize    = 800
	ZipFilesDefaultExpired = 7
)

//ScanLocalFiles 遍历文件
func (z *zipBag) scanLocalFiles() {
	z.mx.Lock()
	defer z.mx.Unlock()

	collectroot := strings.Split(z.config.GetString(ZipCollectRoot), global.SEP)
	for _, v := range collectroot {
		list, err := gfiles.GetAllFiles(v)
		if err != nil {
			z.log.Error("get all files err")
		}
		z.Files = list
	}
}

func (z *zipBag) scanFilesExt() {
	z.mx.Lock()
	defer z.mx.Unlock()

	filesext := strings.Split(z.config.GetString(ZipFilesExt), global.SEP)

	z.Files = filterFilesExt(z.Files, filesext...)
}

func (z *zipBag) scanFilesKeep() {
	z.mx.Lock()
	defer z.mx.Unlock()

	filesfilter := strings.Split(z.config.GetString(ZipFilesKeep), global.SEP)
	z.Files = filterFileKeep(z.Files, filesfilter...)
}

func (z *zipBag) scanFilesLose() {
	z.mx.Lock()
	defer z.mx.Unlock()

	filesfilter := strings.Split(z.config.GetString(ZipFilesLose), global.SEP)
	z.Files = filterFileLose(z.Files, filesfilter...)
}

func (z *zipBag) scanFilesSize() {
	z.mx.Lock()
	defer z.mx.Unlock()

	tempfilse := make(map[string]struct{})

	size := z.config.Get(ZipFilesSize)

	var zipSize int64 = 0
	switch size.(type) {
	case float64:
		zipSize = int64(size.(float64))
	default:
		zipSize = ZipFilesDefaultSize
	}

	for k, v := range z.Files {
		fi, err := os.Stat(k)
		if err != nil {
			z.log.Error("files info err")
			continue
		}

		//文件大小
		if fi.Size() >= zipSize {
			tempfilse[k] = v
		}
	}
	z.Files = tempfilse
}

func (z *zipBag) scanFilesExpired() {
	z.mx.Lock()
	defer z.mx.Unlock()

	tempfilse := make(map[string]struct{})

	expired := z.config.Get(ZipFilesExpired)

	var zipExpired int64 = 0
	switch expired.(type) {
	case float64:
		zipExpired = int64(expired.(float64))

	default:
		zipExpired = ZipFilesDefaultExpired
	}

	for k, v := range z.Files {
		fi, err := os.Stat(k)
		if err != nil {
			continue
		}
		//文件过期
		if time.Now().Unix()-fi.ModTime().Unix() <= zipExpired*24*60*60 {
			tempfilse[k] = v
		}
	}
	z.Files = tempfilse
}

//ExpectFileCache 预期文件vs本地缓存
func (z *zipBag) expectFileCache() {
	z.mx.Lock()
	defer z.mx.Unlock()
	tempfilse := make(map[string]struct{})
	for k, v := range z.Files {
		if value := z.cache.Get(k); value == nil {
			tempfilse[k] = v
			z.cache.Set(k, global.FILENOTREAD)
		}
	}
	z.Files = tempfilse
}

//filterFileExt 获取指定后缀文件
func filterFilesExt(fs map[string]struct{}, exts ...string) map[string]struct{} {
	files := make(map[string]struct{}, 0)
	for k, _ := range fs {
		fileExt := filepath.Ext(k)
		for _, ext := range exts {
			if strings.EqualFold(fileExt, ext) {
				files[k] = struct{}{}
				continue
			}
		}
	}

	return files
}

//filterFileLose 丢弃包含指定字符串的文件
func filterFileLose(fs map[string]struct{}, loses ...string) map[string]struct{} {
	files := make(map[string]struct{}, 0)
	for k, _ := range fs {
		ok := true
		for _, lose := range loses {
			fname := filepath.Base(k)
			if strings.Contains(fname, lose) {
				ok = false
				continue
			}
		}

		if ok {
			files[k] = struct{}{}
		}
	}

	return files
}

//filterFileKeep 保留包含指定字符串的文件
func filterFileKeep(fs map[string]struct{}, keeps ...string) map[string]struct{} {
	files := make(map[string]struct{}, 0)
	for k, _ := range fs {
		for _, keep := range keeps {
			fname := filepath.Base(k)
			if strings.Contains(fname, keep) {
				files[k] = struct{}{}
				continue
			}
		}
	}

	return files
}

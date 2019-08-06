package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strconv"
)

type IConfig interface {
	Get(string) string
	Close()
}

var conf map[string]interface{}

type config struct {
	file string
}

func (this *config) parseConfig() error {

	if len(conf) != 0 {
		return nil
	}

	conf = make(map[string]interface{})
	// 判断json文件是否存在
	if _, err := os.Stat(this.file); os.IsNotExist(err) {
		return err
	}

	// 读取文件内容
	r, err := ioutil.ReadFile(this.file)
	if err != nil {
		return err
	}

	// Unmarshal
	if err := json.Unmarshal(r, &conf); err != nil {
		return err
	}
	return err
}

//Get 获取配置
func (this *config) Get(key string) string {

	if err := this.parseConfig(); err != nil {
		return ""
	}

	value := conf[key]

	switch value.(type) {
	case string:
		return value.(string)
	case bool:
		return strconv.FormatBool(value.(bool))
	case float64:
		return strconv.FormatFloat(value.(float64), 'f', -1, 64)
	default:
		return ""
	}

}

func (this *config) Close() {
	if conf != nil {
		conf = make(map[string]interface{})
	}
}

//New ...
func New(file string) IConfig {
	return &config{
		file: file,
	}
}

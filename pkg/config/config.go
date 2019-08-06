package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"strconv"
)

var configs map[string]interface{}

type config struct {
}

// Get 配置文件
func (this *config) Get(key string) (value interface{}) {
	if configs == nil {
		return
	}
	value = configs[key]
	return
}

// GetString 配置文件
func (this *config) GetString(key string) (value string) {
	v := this.Get(key)

	switch v.(type) {
	case string:
		value = v.(string)
	case bool:
		value = strconv.FormatBool(v.(bool))
	case float64:
		value = strconv.FormatFloat(v.(float64), 'f', -1, 64)
	default:
		value = ""
	}
	return
}

// parseConfigMap json配置文件转成Map
func (this *config) parseConfigMap(file string) (map[string]interface{}, error) {
	if configs != nil {
		return configs, errors.New("already instantiated")
	}

	configs := make(map[string]interface{}, 0)

	// 判断文件是否存在
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return nil, err
	}

	// 读取文件内容
	r, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	// Unmarshal
	if err := json.Unmarshal(r, &configs); err != nil {
		return nil, err
	}

	return configs, err
}

// Close 配置文件
func (this *config) Close() {
	if configs != nil {
		configs = nil
	}
}

//New Config
func New(file string) (c *config) {
	c = &config{}
	configs, _ = c.parseConfigMap(file)
	return
}

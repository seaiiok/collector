package config

import (
	"collector/pkg/interfaces"
	"gcom/gconfig"
)

type config struct {
	configfile string
}

func New(configfile string) interfaces.IConfig {
	gconfig.UpdateConfig(configfile)
	return &config{
		configfile: configfile,
	}
}

func (this *config) Get(k string) string {
	return gconfig.Get(k)
}

func (this *config) Set(k, v string) {
	gconfig.Set(k, v)
}

func (this *config) Close() {
	gconfig.UpdateConfig(this.configfile)
}

package maa

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
)

type Config struct {
	HTTP        *HTTP   `yaml:"http"`
	Asst        []*Asst `yaml:"asst"`
	TaskTimeout string  `yaml:"tasktimeout"`
	StopTimeout string  `yaml:"stoptimeout"`
	UserDir     string  `yaml:"userdir"`
}

type HTTP struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type Asst struct {
	ID           string `yaml:"id"`
	ADBAddress   string `yaml:"adbaddress"`
	ADBPath      string `yaml:"adbpath"`
	ADBType      string `yaml:"adbtype"`
	ResourcePath string `yaml:"resourcepath"`
	UserDir      string `yaml:"userdir"`
	ClientType   string `yaml:"clienttype"`
}

func LoadConfig() (*Config, error) {
	conf := &Config{}
	if err := g.Cfg().MustGet(context.TODO(), "maa").Scan(conf); err != nil {
		return conf, err
	}
	return conf, nil
}

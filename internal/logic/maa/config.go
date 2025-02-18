package maa

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
)

type Config struct {
	HTTP *HTTP	`yaml:"http"`
	Asst *Asst	`yaml:"asst"`
}

type HTTP struct {
	Host string	`yaml:"host"`
	Port int	`yaml:"port"`
}

type Asst struct {
	ADBAddress   string	`yaml:"adb_address"`
	ADBPath      string	`yaml:"adb_path"`
	ADBType      string	`yaml:"adb_type"`
	ResourcePath string	`yaml:"resource_path"`
}

func LoadConfig() (*Config, error) {
	conf := &Config{}
	if err := g.Cfg().MustGet(context.TODO(), "maa").Scan(conf); err != nil {
		return conf, err
	}
	return conf, nil
}
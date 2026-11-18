package maa

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
)

type Config struct {
	HTTP                *HTTP   `yaml:"http"`
	Asst                []*Asst `yaml:"asst"`
	TaskTimeout         string  `yaml:"tasktimeout"`
	StopTimeout         string  `yaml:"stoptimeout"`
	UserDir             string  `yaml:"userdir"`
	ToolboxResourcePath string  `yaml:"toolboxresourcepath"`
}

type HTTP struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type Asst struct {
	ID                  string      `yaml:"id"`
	ADBAddress          string      `yaml:"adbaddress"`
	ADBPath             string      `yaml:"adbpath"`
	ADBType             string      `yaml:"adbtype"`
	ResourcePath        string      `yaml:"resourcepath"`
	UserDir             string      `yaml:"userdir"`
	ClientType          string      `yaml:"clienttype"`
	MuMu                *MuMuConfig `yaml:"mumu"`
	ToolboxResourcePath string      `yaml:"-"`
}

// Index is omitted for ordinary ADB connections; MaaCore derives it from the port.
type MuMuConfig struct {
	Enabled bool   `yaml:"enabled"`
	Path    string `yaml:"path"`
	Index   *int   `yaml:"index"`
}

func LoadConfig() (*Config, error) {
	conf := &Config{}
	if err := g.Cfg().MustGet(context.TODO(), "maa").Scan(conf); err != nil {
		return conf, err
	}
	return conf, nil
}

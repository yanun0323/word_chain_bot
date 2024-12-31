package config

import (
	"sync/atomic"

	"github.com/yanun0323/pkg/config"
	"github.com/yanun0323/pkg/logs"
)

var (
	configuration atomic.Value
)

type Config struct {
	Port               string `json:"port"`
	ChannelID          string `json:"channelID"`
	ChannelSecret      string `json:"channelSecret"`
	ChannelAccessToken string `json:"channelAccessToken"`
}

func LoadConfig() *Config {
	conf, ok := configuration.Load().(*Config)
	if ok {
		return conf
	}

	conf, err := config.InitAndLoad[Config]("config", true, "./config")
	if err != nil {
		logs.Fatalf("load config, err: %+v", err)
	}

	configuration.Store(conf)
	return conf
}

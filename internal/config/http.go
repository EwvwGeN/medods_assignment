package config

import "time"

type HttpConfig struct {
	Host        string  `mapstructure:"host"`
	Port        string  `mapstructure:"port"`
	PingTimeout time.Duration `mapstructure:"ping_timeout"`
}
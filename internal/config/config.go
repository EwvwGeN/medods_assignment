package config

import (
	"fmt"
	p "path"
	"strings"

	"github.com/spf13/viper"
)

var serviceTag string = "auth_service"

type Config struct {
	LogLevel    string      `mapstructure:"log_level"`
	HttpConfig  HttpConfig  `mapstructure:"http"`
}

func LoadConfig(path string) (*Config, error) {
	type ServiceConfig struct {
		Cfg Config `mapstructure:"auth_service"`
	}
	v := viper.NewWithOptions()
	v.AutomaticEnv()
	v.AliasesFirstly(false)
	v.AliasesStepByStep(true)
	if path != "" {
		fileParts := strings.Split(p.Base(path), ".")
		if len(fileParts) < 2 {
			return nil, fmt.Errorf("incorrect config file: %s", path)
		}
		v.SetConfigFile(path)
		v.SetConfigType(fileParts[len(fileParts)-1])
		err := v.ReadInConfig()
		if err != nil {
			return nil, err
		}
	}
	var config ServiceConfig
	keys, err := v.StructKeys(config)
	if err != nil {
		return nil, err
	}
	for _, value := range keys {
		v.RegisterAlias(value, value[len(serviceTag)+1:])
	}
	err = v.Unmarshal(&config)
	if err != nil {
		return nil, err
	}
	return &config.Cfg, nil
}
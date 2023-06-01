package main

import (
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	HttpListen string `mapstructure:"http_listen"`
	Services   struct {
	} `mapstructure:"services"`
}

func ConfigLoad() (*Config, error) {
	viper.SetDefault("http_listen", ":8081")

	// try to read from file
	cfgPath := os.Getenv("CONFIG_PATH")
	if cfgPath == "" {
		cfgPath = "config.yaml"
	}
	viper.SetConfigFile(cfgPath)
	_ = viper.ReadInConfig()

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// try to read from env
	viper.AutomaticEnv()

	conf := &Config{}

	// unmarshal config
	_ = viper.Unmarshal(&conf)

	return conf, nil
}

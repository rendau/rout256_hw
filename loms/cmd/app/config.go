package main

import (
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	GrpcPort string `mapstructure:"grpc_port"`
	HttpPort string `mapstructure:"http_port"`
	Services struct {
	} `mapstructure:"services"`
}

func ConfigLoad() (*Config, error) {
	viper.SetDefault("grpc_port", "8081")
	viper.SetDefault("http_port", "8181")

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

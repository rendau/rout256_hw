package main

import (
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	HttpListen string `mapstructure:"http_listen"`
	Services   struct {
		Loms           ConfigService `mapstructure:"loms"`
		ProductService ConfigService `mapstructure:"product_service"`
	} `mapstructure:"services"`
}

type ConfigService struct {
	Url   string `mapstructure:"url"`
	Token string `mapstructure:"token"`
}

func ConfigLoad() (*Config, error) {
	viper.SetDefault("http_listen", ":8080")
	viper.SetDefault("services.loms.url", "http://localhost:8081")
	viper.SetDefault("services.product_service.url", "")
	viper.SetDefault("services.product_service.token", "")

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

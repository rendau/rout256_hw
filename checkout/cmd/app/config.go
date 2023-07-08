package main

import (
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Debug    bool   `mapstructure:"debug"`
	LogLevel string `mapstructure:"log_level"`
	DbDsn    string `mapstructure:"db_dsn"`
	GrpcPort string `mapstructure:"grpc_port"`
	HttpPort string `mapstructure:"http_port"`
	Services struct {
		Loms           ConfigService `mapstructure:"loms"`
		ProductService ConfigService `mapstructure:"product_service"`
	} `mapstructure:"services"`
	JaegerHostPort string `mapstructure:"jaeger_host_port"`
}

type ConfigService struct {
	Url   string `mapstructure:"url"`
	Token string `mapstructure:"token"`
}

func ConfigLoad() *Config {
	// set default values
	viper.SetDefault("debug", false)
	viper.SetDefault("log_level", "info")
	viper.SetDefault("db_dsn", "postgres://localhost:5432/r256hw_checkout?sslmode=disable")
	viper.SetDefault("grpc_port", "8080")
	viper.SetDefault("http_port", "8180")
	viper.SetDefault("services.loms.url", "localhost:8081")
	viper.SetDefault("services.product_service.url", "")
	viper.SetDefault("services.product_service.token", "")
	viper.SetDefault("jaeger_host_port", "localhost:6831")

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

	return conf
}

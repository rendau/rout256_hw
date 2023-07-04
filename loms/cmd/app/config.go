package main

import (
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	DbDsn    string `mapstructure:"db_dsn"`
	GrpcPort string `mapstructure:"grpc_port"`
	HttpPort string `mapstructure:"http_port"`
	Services struct {
	} `mapstructure:"services"`
	OrderStatusChangeNotifyBrokers []string `mapstructure:"order_status_change_notify_brokers"`
	orderStatusChangeNotifierTopic string   `mapstructure:"order_status_change_notifier_topic"`
}

func ConfigLoad() (*Config, error) {
	viper.SetDefault("db_dsn", "postgres://localhost:5432/r256hw_loms?sslmode=disable")
	viper.SetDefault("grpc_port", "8081")
	viper.SetDefault("http_port", "8181")
	viper.SetDefault("order_status_change_notify_brokers", "")
	viper.SetDefault("order_status_change_notifier_topic", "")

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

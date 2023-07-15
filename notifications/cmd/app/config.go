package main

import (
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Debug                          bool     `mapstructure:"debug"`
	LogLevel                       string   `mapstructure:"log_level"`
	DbDsn                          string   `mapstructure:"db_dsn"`
	GrpcPort                       string   `mapstructure:"grpc_port"`
	HttpPort                       string   `mapstructure:"http_port"`
	KafkaBrokers                   []string `mapstructure:"kafka_brokers"`
	KafkaGroup                     string   `mapstructure:"kafka_group"`
	OrderStatusChangeTopic         string   `mapstructure:"order_status_change_topic"`
	TelegramToken                  string   `mapstructure:"telegram_token"`
	TelegramChatId                 int64    `mapstructure:"telegram_chat_id"`
	OrderStatusChangeEventTemplate string   `mapstructure:"order_status_change_event_template"`
}

func ConfigLoad() *Config {
	// set default values
	viper.SetDefault("debug", false)
	viper.SetDefault("log_level", "info")
	viper.SetDefault("db_dsn", "postgres://localhost:5432/r256hw_notification?sslmode=disable")
	viper.SetDefault("grpc_port", "8082")
	viper.SetDefault("http_port", "8182")
	viper.SetDefault("kafka_brokers", "")
	viper.SetDefault("kafka_group", "")
	viper.SetDefault("order_status_change_topic", "")
	viper.SetDefault("telegram_token", "")
	viper.SetDefault("telegram_chat_id", 0)
	viper.SetDefault("order_status_change_event_template", "Order №%d status changed to %s")

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

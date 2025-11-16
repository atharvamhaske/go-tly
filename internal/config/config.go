package config

import (
	"log"
	"github.com/spf13/viper"
)

type Config struct {
	App struct {
		Name string `mapstructure:"name"`
		Env string  `mapstructure:"env"`
		Port string `mapstructure:"port"`
	} `mapstructure:"app"`

	MongoDB struct {
		URI string `mapstructure:"uri"`
		Database string `mapstructure:"database"`
	} `mapstructure:"mongodb"`

	Redis struct {
		Addr string `mapstructure:"addr"`
		Password string `mapstructure:"password"`
		DB string `mapstructure:"db"`
	} `mapstructure:"redis"`

	ClickHouse struct {
		DSN string `mapstructure:"dsn"`
	} `mapstructure:"clickhouse"`

	GRPC struct {
		KGSServerAddr string `mapstructure:"kgs_server_addr"`
	} `mapstructure:"grpc"`
}

func LoadConfig() (*Config, error) {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.AddConfigPath("./config")

	v.AutomaticEnv() // here env overrides YAML

	if err := v.ReadInConfig();
	err != nil {
		log.Printf("No config file found, using env only: %v", err)
	}

	cfg := &Config{}
		if err := v.Unmarshal(cfg);
        err != nil {
           return nil, err
        }
           return cfg, nil
}

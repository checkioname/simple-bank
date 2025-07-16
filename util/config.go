package util

import (
	"github.com/spf13/viper"
	"time"
)

const (
	connStr = "postgres://root:secret@localhost/simple_bank?sslmode=disable"
	addr    = "0.0.0.0:8080"
)

type Config struct {
	ConnStr             string        `mapstructure:"CONN_STR"`
	Addr                string        `mapstructure:"ADDR"`
	TokenSymmetricKey   string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
}

func LoadConfig() (config Config, err error) {
	viper.SetConfigName("app") // name of config file (without extension)
	viper.SetConfigType("env") // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}

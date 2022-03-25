package conf

import (
	"github.com/spf13/viper"
)

// Config is app config
type Config struct {
	Prod      bool   `yaml:"Prod"`
	Addr      string `yaml:"Addr"`
	RedisHost string `yaml:"RedisHost"`
}

// AppConfig everything we need for app config
var AppConfig Config

// Init init config
func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./conf")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	err := viper.Unmarshal(&AppConfig)
	if err != nil {
		panic(err)
	}
}

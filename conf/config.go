package conf

import (
	"github.com/spf13/viper"
)

// HystrixModelConf is hystrix config
type HystrixModelConf struct {
	Timeout                int
	MaxConcurrentRequests  int
	RequestVolumeThreshold int
	ErrorPercentThreshold  int
}

// HystrixConf is hystrix config
type HystrixConf struct {
	Degrade HystrixModelConf
}

// DashBoradConf is dashboard config
type DashBoradConf struct {
}

// Config is app config
type Config struct {
	ENV       string      `yaml:"ENV"`
	Addr      string      `yaml:"Addr"`
	RedisHost string      `yaml:"RedisHost"`
	Hystrix   HystrixConf `yaml:"Hystrix"`
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

// IsProd return true if app is prod
func IsProd() bool {
	return AppConfig.ENV == "prod"
}

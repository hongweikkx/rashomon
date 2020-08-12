package conf

import (
	"github.com/spf13/viper"
)

type HttpServerConf struct {
	Addr string
}

type GrpcServerConf struct {
	Addr string
}

type HystrixModelConf struct {
	Timeout int
	MaxConcurrentRequests int
	RequestVolumeThreshold int
	ErrorPercentThreshold int
}

type HystrixConf struct {
	Degrade HystrixModelConf
	Fuse HystrixModelConf
}

type JWTConf struct {
	Enable bool
}

type ETCDConf struct {
	Enable bool
	EndPoints []string
	DailTimeout int
}

type LoadBalanceConf struct {
	Algorithm int
}

type Config struct {
	HttpServer HttpServerConf `yaml:"HttpServer"`
	GrpcServer GrpcServerConf `yaml:"GrpcServer"`
	Hystrix HystrixConf `yaml:"Hystrix"`
	JWT JWTConf `yaml:"JWT"`
	ETCD ETCDConf `yaml:"ETCD"`
	LoadBalance LoadBalanceConf `yaml:"LoadBalance"`
}

var AppConfig Config

func InitConf() error{
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./conf")
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	err := viper.Unmarshal(&AppConfig)
	if err != nil {
		return err
	}
	return nil
}
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

type Config struct {
	HttpServer HttpServerConf `yaml:"httpServer"`
	GrpcServer GrpcServerConf `yaml:"grpcServer"`
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
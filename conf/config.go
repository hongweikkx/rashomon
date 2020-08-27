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
	User string
	Password string
	EndPoints []string
	DailTimeout int
	WatchPrix string
}

//Algorithm: 2   # WeightedRoundRobin = 1 / RoundRobin = 2 / consistentHashing = 3  | default = 2
type LoadBalanceConf struct {
	Algorithm int
}


type ProxyConf struct {
	HttpServer HttpServerConf `yaml:"HttpServer"`
	GrpcServer GrpcServerConf `yaml:"GrpcServer"`
}

type DiscoveryConf struct {
	Service string
	ETCD ETCDConf `yaml:"ETCD"`
}

type Config struct {
	Proxy ProxyConf `yaml:"Proxy"`
	Discovery DiscoveryConf `yaml:"Discovery"`
	Hystrix HystrixConf `yaml:"Hystrix"`
	JWT JWTConf `yaml:"JWT"`
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
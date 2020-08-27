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


type ETCDConf struct {
	User string
	Password string
	EndPoints []string
	DailTimeout int
	WatchPrix string
}

//Algorithm: 2   # WeightedRoundRobin = 1 / RoundRobin = 2 / consistentHashing = 3  | default = 2

type ProxyConf struct {
	HttpServer HttpServerConf `yaml:"HttpServer"`
	GrpcServer GrpcServerConf `yaml:"GrpcServer"`
}

type StorageConf struct {
	Service string
	ETCD ETCDConf `yaml:"ETCD"`
}

type DashBoradConf struct {
	Addr string
}

type Config struct {
	ENV   string    `yaml:"ENV"`
	Proxy ProxyConf `yaml:"Proxy"`
	DashBoard DashBoradConf `yaml:"DashBoard"`
	Storage StorageConf `yaml:"Storage"`
	Hystrix HystrixConf `yaml:"Hystrix"`
}

var AppConfig Config

func Init() error{
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

func IsProd() bool{
	return AppConfig.ENV == "prod"
}
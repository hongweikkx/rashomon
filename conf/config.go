package conf

import (
	"github.com/spf13/viper"
)

// HTTPServerConf is http server config
type HTTPServerConf struct {
	Addr string
}

// GrpcServerConf is grpc server config
type GrpcServerConf struct {
	Addr string
}

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
	Fuse    HystrixModelConf
}

// ETCDConf is etcd config
type ETCDConf struct {
	User        string
	Password    string
	EndPoints   []string
	DailTimeout int
	WatchPrix   string
}

//Algorithm: 2   # WeightedRoundRobin = 1 / RoundRobin = 2 / consistentHashing = 3  | default = 2

// ProxyConf is proxy config
type ProxyConf struct {
	HTTPServer HTTPServerConf `yaml:"HTTPServer"`
	GrpcServer GrpcServerConf `yaml:"GrpcServer"`
}

// StorageConf is storage config
type StorageConf struct {
	Service string
	ETCD    ETCDConf `yaml:"ETCD"`
}

// DashBoradConf is dashboard config
type DashBoradConf struct {
	Addr string
}

// Config is app config
type Config struct {
	ENV       string        `yaml:"ENV"`
	Proxy     ProxyConf     `yaml:"Proxy"`
	DashBoard DashBoradConf `yaml:"DashBoard"`
	Storage   StorageConf   `yaml:"Storage"`
	Hystrix   HystrixConf   `yaml:"Hystrix"`
}

// AppConfig everything we need for app config
var AppConfig Config

// Init init config
func Init() error {
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

// IsProd return true if app is prod
func IsProd() bool {
	return AppConfig.ENV == "prod"
}

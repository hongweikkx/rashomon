package conf

import (
	"github.com/spf13/viper"
)

var AppConfig Config

type ServiceConf struct {
	AppMode  string `yaml:"AppMode"`
	AppName  string `yaml:"AppName"`
	HttpPort string `yaml:"HttpPort"`
}

type MysqlConf struct {
	Db         string `yaml:"Db"`
	DbHost     string `yaml:"DbHost"`
	DbPort     string `yaml:"DbPort"`
	DbUser     string `yaml:"DbUser"`
	DbPassWord string `yaml:"DbPassWord"`
	DbName     string `yaml:"DbName"`
}

type RedisConf struct {
	RedisDb       string `yaml:"RedisDb"`
	RedisAddr     string `yaml:"RedisAddr"`
	RedisPw       string `yaml:"RedisPw"`
	RedisDbName   string `yaml:"RedisDbName"`
	RedisWarnTime int64  `yaml:"RedisWarnTime"`
}

type HttpConf struct {
	ReadTimeout    int64 `yaml:"ReadTimeout"`
	WriteTimeout   int64 `yaml:"WriteTimeout"`
	MaxHeaderBytes uint  `yaml:"MaxHeaderBytes"`
}

type OSS struct {
	EndPoint        string `yaml:"EndPoint"`
	AccessKeyId     string `yaml:"AccessKeyId"`
	AccessKeySecret string `yaml:"AccessKeySecret"`
	SecurityToken   string `yaml:"SecurityToken"`
	Bucket          string `yaml:"Bucket"`
	Address         string `yaml:"Address"`
}

type LogConf struct {
	Level      string `json:"level"`       // Level 最低日志等级，DEBUG<INFO<WARN<ERROR<FATAL 例如：info-->收集info等级以上的日志
	FileName   string `json:"file_name"`   // FileName 日志文件位置
	MaxSize    int    `json:"max_size"`    // MaxSize 进行切割之前，日志文件的最大大小(MB为单位)，默认为100MB
	MaxAge     int    `json:"max_age"`     // MaxAge 是根据文件名中编码的时间戳保留旧日志文件的最大天数。
	MaxBackups int    `json:"max_backups"` // MaxBackups 是要保留的旧日志文件的最大数量。默认是保留所有旧的日志文件（尽管 MaxAge 可能仍会导致它们被删除。）
}

type Config struct {
	Service ServiceConf `yaml:"service"`
	Mysql   MysqlConf   `yaml:"mysql"`
	Redis   RedisConf   `yaml:"redis"`
	Http    HttpConf    `yaml:"http"`
	Log     LogConf     `yaml:"logger"`
	Oss     OSS         `yaml:"OSS"`
}

// Init init config
func Init() {
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

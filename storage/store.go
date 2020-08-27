package storage

import (
	"errors"
	"github.com/hongweikkx/rashomon/conf"
	"github.com/hongweikkx/rashomon/storage/etcd"
)

// 实现此接口才能作为discovery
type Storeage interface {
	WatchWorkers(string)
	Add(keyByte []byte, serviceByte []byte)
	Delete(keyByte []byte)
	Stop()
}

func StartStorage() (Storeage, error){
	// 应该是每个proxy有个etcd
	switch conf.AppConfig.Discovery.Service {
	case "etcd":
		master, err := etcd.New(etcd.Auth{User: conf.AppConfig.Discovery.ETCD.User, Password: conf.AppConfig.Discovery.ETCD.Password})
		return master, err
	default:
		return nil, errors.New("disovery do not support:" + conf.AppConfig.Discovery.Service)
	}
}


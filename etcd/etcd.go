package etcd

import (
	"context"
	"errors"
	"github.com/hongweikkx/rashomon/conf"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/etcdserver/api/v3rpc/rpctypes"
	"time"
)

func New() (*clientv3.Client, error){
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   conf.AppConfig.ETCD.EndPoints,
		DialTimeout:  time.Duration(conf.AppConfig.ETCD.DailTimeout) * time.Second,
	})
	if err != nil {
		return nil, err
	}
	return cli, nil
}

func Put(cli *clientv3.Client, key, value string) error{
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(conf.AppConfig.ETCD.DailTimeout) * time.Second)
	_, err := cli.Put(ctx, key, value)
	cancel()
	if err != nil {
		switch err {
		case context.Canceled:
			return errors.New("ctx is canceled by another routine:" + err.Error())
		case context.DeadlineExceeded:
			return errors.New("ctx is attached with a deadline is exceeded: " + err.Error())
		case rpctypes.ErrEmptyKey:
			return errors.New("client-side error: " + err.Error())
		default:
			return errors.New("bad cluster endpoints, which are not etcd servers: " + err.Error())
		}
	}
	return nil
}
func Get(cli *clientv3.Client, key, value string) (*clientv3.GetResponse, error){
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(conf.AppConfig.ETCD.DailTimeout) * time.Second)
	resp, err := cli.Get(ctx, key)
	cancel()
	if err != nil {
		switch err {
		case context.Canceled:
			return nil, errors.New("ctx is canceled by another routine:" + err.Error())
		case context.DeadlineExceeded:
			return nil, errors.New("ctx is attached with a deadline is exceeded: " + err.Error())
		case rpctypes.ErrEmptyKey:
			return nil, errors.New("client-side error: " + err.Error())
		default:
			return nil, errors.New("bad cluster endpoints, which are not etcd servers: " + err.Error())
		}
	}
	return resp, nil
}

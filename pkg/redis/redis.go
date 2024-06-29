package redis

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"rashomon/conf"
	"rashomon/consts"
	"rashomon/pkg/logger"
	"runtime"
	"sync"
	"time"
)

var (
	once  sync.Once //确保全局只实例化一次redis
	Redis *Client
)

type Client struct {
	*redis.Client
}

func Init() {
	once.Do(func() {
		var initRedisErr error
		Redis, initRedisErr = newRedisClient()
		if initRedisErr != nil {
			panic(initRedisErr)
		}
	})
}

func newRedisClient() (*Client, error) {
	// 初始化自定的 RedisClient 实例
	rds := &Client{}
	// 使用默认的 context

	// 使用 redis 库里的 NewClient 初始化连接
	rds.Client = redis.NewClient(&redis.Options{
		Addr:         conf.AppConfig.Redis.RedisAddr,
		Password:     conf.AppConfig.Redis.RedisPw,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 2 * time.Second,
		PoolSize:     runtime.NumCPU() / 2,
		MinIdleConns: 1,
	})

	// 测试一下连接
	if err := rds.Ping(context.Background()); err != nil {

	}
	return rds, nil
}

// Ping 用以测试 redis 连接是否正常
func (rds *Client) Ping(ctx context.Context) error {
	defer rds.Metric(time.Now())
	_, err := rds.Client.Ping(ctx).Result()
	return err
}

// Metric redis timeout metric
func (rds *Client) Metric(start time.Time) {
	duration := time.Since(start)
	if duration.Milliseconds() >= conf.AppConfig.Redis.RedisWarnTime {
		pc, file, line, _ := runtime.Caller(1)
		logger.Logger.Warn(consts.METRIC_REDIS, zap.Int64("duration", duration.Milliseconds()),
			zap.Any("function", runtime.FuncForPC(pc).Name()), zap.Any("file", file), zap.Any("line", line))
	}
}

// Set 存储 key 对应的 value，且设置 expiration 过期时间
func (rds *Client) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) (string, error) {
	defer rds.Metric(time.Now())
	res, err := rds.Client.Set(ctx, key, value, expiration).Result()
	rds.ErrLog(err)
	return res, err
}

// GetExists 获取 key 对应的 value
func (rds *Client) GetExists(ctx context.Context, key string) (string, error) {
	defer rds.Metric(time.Now())
	res, err := rds.Client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return "", nil
	}
	rds.ErrLog(err, redis.Nil)
	return res, err
}

// Get 不可以忽略 redis.Nil
func (rds *Client) Get(ctx context.Context, key string) (string, error) {
	defer rds.Metric(time.Now())
	res, err := rds.Client.Get(ctx, key).Result()
	rds.ErrLog(err, redis.Nil)
	return res, err
}

// Has 判断一个 key 是否存在
func (rds *Client) Has(ctx context.Context, key string) (bool, error) {
	defer rds.Metric(time.Now())
	res, err := rds.Client.Exists(ctx, key).Result()
	rds.ErrLog(err)
	if err != nil {
		return false, err
	}
	return res != 0, nil
}

// Del 删除存储在 redis 里的数据
func (rds *Client) Del(ctx context.Context, keys ...string) (int64, error) {
	defer rds.Metric(time.Now())
	res, err := rds.Client.Del(ctx, keys...).Result()
	rds.ErrLog(err)
	return res, err
}

// Lock 分布式锁-加锁
func (rds *Client) Lock(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	defer rds.Metric(time.Now())
	res, err := rds.Client.SetNX(ctx, key, value, expiration).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			err = nil
		}
		return false, err
	}
	rds.ErrLog(err)
	return res, err
}

// KeepLock 分布式锁-保持锁的过期时间
func (rds *Client) KeepLock(ctx context.Context, key string, expiration time.Duration) (bool, error) {
	defer rds.Metric(time.Now())
	res, err := rds.Client.Expire(ctx, key, expiration).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			err = nil
		}
		return false, err
	}
	rds.ErrLog(err)
	return res, err
}

// Unlock 分布式锁-解锁
func (rds *Client) Unlock(ctx context.Context, key string) (int64, error) {
	defer rds.Metric(time.Now())
	res, err := rds.Client.Del(ctx, key).Result()
	rds.ErrLog(err)
	return res, err
}

// ================================== sets =================================

func (rds *Client) SAdd(ctx context.Context, key string, members ...interface{}) (int64, error) {
	defer rds.Metric(time.Now())
	res, err := rds.Client.SAdd(ctx, key, members...).Result()
	rds.ErrLog(err)
	return res, err
}

func (rds *Client) SCard(ctx context.Context, key string) (int64, error) {
	defer rds.Metric(time.Now())
	res, err := rds.Client.SCard(ctx, key).Result()
	rds.ErrLog(err)
	return res, err
}

func (rds *Client) SRem(ctx context.Context, key string, members ...interface{}) (int64, error) {
	defer rds.Metric(time.Now())
	res, err := rds.Client.SRem(ctx, key, members...).Result()
	rds.ErrLog(err)
	return res, err
}

func (rds *Client) SIsMember(ctx context.Context, key string, member interface{}) (bool, error) {
	defer rds.Metric(time.Now())
	res, err := rds.Client.SIsMember(ctx, key, member).Result()
	rds.ErrLog(err)
	return res, err
}

// ================================== lists =================================

func (rds *Client) RPush(ctx context.Context, key string, values ...interface{}) (int64, error) {
	defer rds.Metric(time.Now())
	res, err := rds.Client.RPush(ctx, key, values...).Result()
	rds.ErrLog(err)
	return res, err
}

func (rds *Client) LPush(ctx context.Context, key string, values ...interface{}) (int64, error) {
	defer rds.Metric(time.Now())
	res, err := rds.Client.LPush(ctx, key, values...).Result()
	rds.ErrLog(err)
	return res, err
}

func (rds *Client) RPop(ctx context.Context, key string) (string, error) {
	defer rds.Metric(time.Now())
	res, err := rds.Client.RPop(ctx, key).Result()
	rds.ErrLog(err)
	return res, err
}

func (rds *Client) LLen(ctx context.Context, key string) (int64, error) {
	defer rds.Metric(time.Now())
	res, err := rds.Client.LLen(ctx, key).Result()
	rds.ErrLog(err)
	return res, err
}

func (rds *Client) ErrLog(err error, ignoreErr ...error) {
	//if err != nil && !helpers.IsErrListContains(ignoreErr, err) {
	//	logger.LogIf(err)
	//}
}

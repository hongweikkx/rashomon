package user

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"rashomon/dao"
	"rashomon/model"
	"rashomon/pkg/logger"
	redisDB "rashomon/pkg/redis"
	"time"
)

type Service struct{}

func NewUserService() Service {
	return Service{}
}

func (us Service) GetUserByName(c context.Context, userName string) (*model.User, error) {
	userRedis, err := redisDB.Redis.Get(c, userName)
	if errors.Is(err, redis.Nil) {
		userDb, err := dao.NewUserDao(c).GetUserByName(userName)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("can not find the user")
		}
		if err != nil {
			return nil, err
		}
		us.CacheUser(c, userDb)
		return userDb, nil
	}
	if err != nil {
		return nil, err
	}
	user := model.User{}
	err = json.Unmarshal([]byte(userRedis), &user)
	return &user, err
}

func (us Service) CacheUser(c context.Context, user *model.User) {
	marshalUser, err := json.Marshal(*user)
	if err != nil {
		logger.Error(c, "cache user", zap.Error(err))
		return
	}
	_, err = redisDB.Redis.Set(c, user.Name, string(marshalUser), time.Minute)
	if err != nil {
		logger.Error(c, "cache user", zap.Error(err))
	}
	return
}

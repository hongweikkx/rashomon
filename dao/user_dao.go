package dao

import (
	"context"
	"rashomon/model"
	"rashomon/pkg/mysql"
)

type UserDao struct {
	TableName string
}

// NewUserDao 实例化Dao对象
func NewUserDao(ctx context.Context) *UserDao {
	return &UserDao{"user"}
}

// GetUserById 根据 id 获取用户
func (dao *UserDao) GetUserById(id uint64) (user *model.User, err error) {
	err = mysql.DB.Model(&model.User{}).Where("id=?", id).First(&user).Error
	return
}

// GetUserByName 根据 name 获取用户
func (dao *UserDao) GetUserByName(name string) (user *model.User, err error) {
	err = mysql.DB.Model(&model.User{}).Where("name=?", name).First(&user).Error
	return
}

package dao

import (
	"context"
	"gorm.io/gorm"
	"rashomon/model"
	"rashomon/pkg/mysql"
)

type UserDao struct {
	DB        *gorm.DB
	TableName string
}

// NewUserDao 实例化Dao对象
func NewUserDao(c context.Context) *UserDao {
	return &UserDao{mysql.NewWithContext(c), "user"}
}

// GetUserById 根据 id 获取用户
func (dao *UserDao) GetUserById(id uint64) (user *model.User, err error) {
	err = dao.DB.Model(&model.User{}).Where("id=?", id).First(&user).Error
	return
}

// GetUserByName 根据 name 获取用户
func (dao *UserDao) GetUserByName(name string) (user *model.User, err error) {
	err = dao.DB.Model(&model.User{}).Where("name=?", name).First(&user).Error
	return
}

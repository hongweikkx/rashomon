package model

import (
	"github.com/spf13/cast"
	"rashomon/pkg/localtime"
)

type BaseModel struct {
	ID uint64 `gorm:"column:id;primaryKey;autoIncrement;" json:"id"`
}

// CommonTimestampsField 时间戳
type CommonTimestampsField struct {
	CreatedAt *localtime.LocalTime `gorm:"column:created_at;index;" json:"created_at"`
	UpdatedAt *localtime.LocalTime `gorm:"column:updated_at;index;" json:"updated_at"`
}

// GetStringID 获取 ID 的字符串格式
func (a BaseModel) GetStringID() string {
	return cast.ToString(a.ID)
}

// DeletedAtField 软删除时间
type DeletedAtField struct {
	DeletedAt *localtime.LocalTime `gorm:"column:deleted_at;index;" json:"deleted_at" sql:"DEFAULT:NULL"`
}

package localtime

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type LocalTime time.Time

// MarshalJSON 返序列化时候将时间字符串改成标准格式
func (t *LocalTime) MarshalJSON() ([]byte, error) {
	tTime := time.Time(*t)
	return []byte(fmt.Sprintf("\"%v\"", tTime.Local().Format("2006-01-02 15:04:05"))), nil
}

// UnmarshalJSON 反序列化时候将将字符串改成标准格式
func (t *LocalTime) UnmarshalJSON(data []byte) error {
	tmp, err := time.ParseInLocation("\"2006-01-02 15:04:05\"", string(data), time.Now().Location())
	if err != nil {
		return err
	}
	*t = LocalTime(tmp)
	return nil
}

func (t LocalTime) ToTime() time.Time {
	return time.Time(t)
}

// Value 时间入库时候转换成timestamp
func (t LocalTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	tlt := time.Time(t)
	//判断给定时间是否和默认零时间的时间戳相同
	if tlt.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return tlt, nil
}

// Scan 读取数据时候转换成 LocalTime
func (t *LocalTime) Scan(v interface{}) error {
	if value, ok := v.(time.Time); ok {
		*t = LocalTime(value)
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}

package store

import (
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Mysql struct {
	db *gorm.DB
}

func (m Mysql) Get(code, key string) (interface{}, error) {
	//TODO implement me
	panic("implement me")
}

func (m Mysql) GetString(code, key string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (m Mysql) GetLast(code string) (string, interface{}) {
	//TODO implement me
	panic("implement me")
}

func (m Mysql) Set(code, key string, value interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (m Mysql) Delete(code, key string) {
	//TODO implement me
	panic("implement me")
}

func NewMysql(dsn string) *Mysql {
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       dsn,   // DSN data source name
		DefaultStringSize:         256,   // string 类型字段的默认长度
		SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
	}), &gorm.Config{})
	if err != nil {
		logrus.Fatalf("gorm.Open err: %v", err)
	}
	return &Mysql{
		db: db,
	}
}

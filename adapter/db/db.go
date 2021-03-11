/*
 * @Descripttion:
 * @version:
 * @Author: Cmpeax
 * @Date: 2019-12-18 11:43:25
 * @LastEditors  : Cmpeax
 * @LastEditTime : 2019-12-18 11:43:56
 */
package db

import (
	"errors"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"xorm.io/xorm"
	"xorm.io/xorm/log"
)

// Caches     *redis.Pool
// "github.com/gomodule/redigo/redis"
type DB struct {
	db *xorm.Engine
}

func NewDB() *DB {
	return &DB{}
}

func (d *DB) Get() *xorm.Engine {
	return d.db
}

// 创建表 如果没有被创建
//func (t *DB) CreateTableIfNotExist(tableName string, tableStruct ITable) error {
//	isExist, err := t.Get().IsTableExist(tableName)
//	if err != nil {
//		return err
//	}
//	if isExist == false {
//		tableStruct.SetTableName(tableName)
//		t.Get().CreateTables(tableStruct)
//	}
//	return nil
//}

func (s *DB) Install(param ...interface{}) error {

	model, ok := param[0].(string)
	if !ok {
		return errors.New("第一个参数 model 丢失!")
	}
	urls, ok := param[1].(string)
	if !ok {
		return errors.New("第二个参数 urls 丢失!")
	}

	engine, err := xorm.NewEngine(model, urls)
	if err != nil {
		return err
	}

	// 生成日志
	f, err := os.Create("sql.log")
	if err != nil {
		return err
	}
	engine.SetLogger(log.NewSimpleLogger(f))
	engine.ShowSQL(true)
	engine.SetMaxIdleConns(10)
	engine.SetMaxOpenConns(100)
	engine.SetConnMaxLifetime(100 * time.Second)

	// 测试是否能链接到数据库
	if err = engine.Ping(); err != nil {
		return err
	}

	s.db = engine
	return nil

}

func (s *DB) Start(param ...interface{}) error {
	// 等待发送函数

	return nil
}

func (s *DB) Uninstall(param ...interface{}) error {
	s.db.Close()
	return nil
}

// func (m *DataBase) Set_Cache(addr string) {
// 	m.Caches = &redis.Pool{
// 		MaxIdle:     3,
// 		IdleTimeout: 240 * time.Second,
// 		// Dial or DialContext must be set. When both are set, DialContext takes precedence over Dial.
// 		Dial: func() (redis.Conn, error) {
// 			c, err := redis.Dial("tcp", addr)
// 			if err != nil {
// 				return nil, err
// 			}
// 			if _, err := c.Do("AUTH", "57RXjkhCZEzx1z70"); err != nil {
// 				c.Close()
// 				return nil, err
// 			}
// 			if _, err := c.Do("SELECT", "1"); err != nil {
// 				c.Close()
// 				return nil, err
// 			}
// 			return c, nil

// 		},
// 	}
// }

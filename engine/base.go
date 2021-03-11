/**
 * @Author: cmpeax Tang
 * @Date: 2021/3/10 15:57
 */

package engine

import (
	"esd-router-preview/adapter/cache"
	"esd-router-preview/core/eventbus"
	"io/ioutil"
	"sync"
)

// ESD引擎
type TypeESDEngine struct {
	configEventBus string

	orm      IORM         // 我需要数据库组件
	cache    cache.ICache // 我需要缓存组件
	eventBus IEventBus    // 消息总线
}

var instance *TypeESDEngine
var once sync.Once

// 导入 描述通道型 JSON文件.
func GetInstance() *TypeESDEngine {
	once.Do(func() {
		instance = &TypeESDEngine{}
	})
	return instance
}

func (t *TypeESDEngine) GetEventBus() IEventBus {
	return t.eventBus
}

// 读取配置文件
func (t *TypeESDEngine) LoadConfig(filePath string) {
	config, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err.Error())
	}
	t.configEventBus = string(config)
	// 转换 成json

}

// 执行预加载
func (t *TypeESDEngine) Init() error {

	// 初始化 装入 基础组件
	// 消息总线

	tEventBus, err := eventbus.NewEventBus(t.configEventBus)
	if err != nil {
		return err
	}
	err = tEventBus.Install()
	if err != nil {
		return err
	}
	// 挂载事件:
	//tEventBus.Mount(func(data transmit.TypeMessageChannel) []transmit.TypeMessageChannel {
	//	return []transmit.TypeMessageChannel{}
	//})

	t.eventBus = tEventBus

	return nil
}

// 跑起来
func (t *TypeESDEngine) Start() {
	t.eventBus.Start()
}

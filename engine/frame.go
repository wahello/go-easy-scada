package engine

import (
	"esd-router-preview/adapter"
	"esd-router-preview/adapter/transmit"
)

// 面向接口的相关功能定义

// orm框架 (行为规范暂时比较困难)
type IORM interface {
	adapter.IBase
	Get() interface{}
}

// 数据源 （强调复用）
type ISource interface {
	adapter.IBase
	GetReader() chan interface{}
	GetWriter() chan interface{}
}

// 定时器组件
type ITicker interface {
	adapter.IBase
}

// 消息总线 发送接收对
type IEventBus interface {
	IBase
	Mount(func(transmit.TypeMessageChannel) []transmit.TypeMessageChannel)
	Send(transmit.TypeMessageChannel)
}

// 定义了接口的生命周期
type IBase interface {
	Install(...interface{}) error
	Start(...interface{})
	Uninstall(...interface{}) error
}

/**
 * @Author: cmpeax Tang
 * @Date: 2021/3/10 21:51
 */

package eventbus

import (
	"errors"
	"esd-router-preview/adapter/transmit"
	_ "esd-router-preview/adapter/transmit/natsmq"
	_ "esd-router-preview/adapter/transmit/webpost"
	"esd-router-preview/core/logout"
	"fmt"
)

type TypeEventBusService struct {
	mounted     func(transmit.TypeMessageChannel) []transmit.TypeMessageChannel
	handleWrite chan transmit.TypeMessageChannel
	config      TypeDefineHook
	inputMap    map[string][]transmit.IInputSource
	outputMap   map[string][]transmit.IOutputSource
}

func (t *TypeEventBusService) Send(data transmit.TypeMessageChannel) {
	t.handleWrite <- data
}

func (t *TypeEventBusService) Install(param ...interface{}) error {
	return t.init()
}

// 消息总线初始化
func (t *TypeEventBusService) init() error {
	mapISource := make(map[string]transmit.ISource, 0)
	//
	//// 开始读取插件并导入.
	for k, v := range t.config.Plugin {
		sourceObject, err := transmit.New(v.Module, v.Param, k)
		if err != nil {
			return errors.New(fmt.Sprintf("检测plugin参数时,发现 %s 出现错误: %s", k, err.Error()))
		}

		sourceObject.Start()

		mapISource[k] = sourceObject
	}

	// mapISource: 键值对 键：插件名，对：对应对象
	// 开始读取事件并导入.
	mapMsgInput := make(map[string][]transmit.IInputSource, 0)
	mapMsgOutput := make(map[string][]transmit.IOutputSource, 0)

	// 按事件类型进行遍历.
	for k, v := range t.config.Message {
		tempPubSlice := []transmit.IInputSource{}
		tempSubSlice := []transmit.IOutputSource{}

		// pub input
		for _, item := range v.Pub {
			if _, ok := mapISource[item]; !ok {
				return errors.New("找不到插件:" + item)
			}
			if !mapISource[item].IsInputer() {
				return errors.New("插件:" + item + "不支持或者不实现接收方法!")
			}
			tempPubSlice = append(tempPubSlice, mapISource[item])
		}

		// sub output
		for _, item := range v.Sub {
			if _, ok := mapISource[item]; !ok {
				return errors.New("找不到插件:" + item)
			}
			if !mapISource[item].IsOutputer() {
				return errors.New("插件:" + item + "不支持或者不实现接收方法!")
			}
			tempSubSlice = append(tempSubSlice, mapISource[item])
		}
		mapMsgInput[k] = tempPubSlice
		mapMsgOutput[k] = tempSubSlice
	}
	//// fmt.Println(t.config)
	t.inputMap = mapMsgInput
	t.outputMap = mapMsgOutput
	t.handleWrite = make(chan transmit.TypeMessageChannel, 200)

	return nil
}

func (t *TypeEventBusService) Mount(cb func(transmit.TypeMessageChannel) []transmit.TypeMessageChannel) {
	t.mounted = cb
}

// 回调
func (t *TypeEventBusService) callback(pubChan chan transmit.TypeMessageChannel, data transmit.TypeMessageChannel) {
	if t.mounted != nil {
		result := t.mounted(data)
		for _, v := range result {
			go func(tv transmit.TypeMessageChannel) {
				pubChan <- tv
			}(v)
		}
	}

}

func (t *TypeEventBusService) Start(param ...interface{}) {
	logout.L.PrintLog(logout.L.Info(), "消息总线启动!")
	if t.mounted == nil {
		logout.L.PrintLog(logout.L.Info(), "发现你并没有挂载回调消息在 消息总线 上.!")
	}

	// 接收 协程
	go func(maps map[string][]transmit.IInputSource) {
		for k, _ := range maps {
			go func(key string) {
				pubSlice := maps[key]
				for i, _ := range pubSlice {
					go func(index int, tKey string) {
						for {
							// 从订阅者 pub 获取信息.
							tdata := <-pubSlice[index].GetReader()
							go t.callback(t.handleWrite, tdata)
							t.handleWrite <- transmit.TypeMessageChannel{
								Data:    tdata.Data,
								Message: tKey,
								Param:   tdata.Param,
							}

						}
					}(i, key)
				}
			}(k)
		}
	}(t.inputMap)

	// 发送 协程
	go func(reader chan transmit.TypeMessageChannel, maps map[string][]transmit.IOutputSource) {
		for {
			data := <-reader
			msg := data.Message
			if _, ok := maps[msg]; !ok {
				logout.L.PrintLog(logout.L.Info(), "找不到对应的消息:", msg)
				continue
			}
			for _, v := range maps[msg] {
				go func(subobj transmit.IOutputSource) {
					subobj.GetWriter() <- data
				}(v)
			}
		}
	}(t.handleWrite, t.outputMap)

	// Sub的多路复用
}

func (t *TypeEventBusService) Uninstall(param ...interface{}) error {
	return nil
}

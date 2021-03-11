/**
 * @Author: cmpeax Tang
 * @Date: 2021/3/12 3:32
 */

package natsmq

import (
	"errors"
	"esd-router-preview/adapter/transmit"
	"esd-router-preview/utilio"
)

func NewNatsInput() transmit.ISource {
	return &NatsInput{}
}

type NatsInput struct {
	addr  string
	token string
	topic string

	handleRead chan transmit.TypeMessageChannel
}

func (s *NatsInput) GetReader() chan transmit.TypeMessageChannel {
	return s.handleRead
}

func (s *NatsInput) GetWriter() chan transmit.TypeMessageChannel {
	return nil
}

func (s *NatsInput) IsInputer() bool {
	return true
}

func (s *NatsInput) IsOutputer() bool {
	return false
}

func (n *NatsInput) SetToken(token string) {
	n.token = token
}

func (s *NatsInput) Start(param ...interface{}) {
	utilio.AdapterInfoLog(s.token, pluginName, "启动服务")
}

func (s *NatsInput) Install(param string) error {
	s.handleRead = make(chan transmit.TypeMessageChannel, 100)

	tempParams := map[string]*string{
		"addr":  &s.addr,
		"topic": &s.topic,
	}

	if err := utilio.ValidateConnStr(&tempParams, param); err != nil {
		return errors.New(utilio.GetAdapterLog(s.token, pluginName, err.Error()))
	}

	// 等待发送函数
	return mqInstance.NewInput(s.token, s.addr, s.topic, func(data interface{}) {
		s.handleRead <- transmit.TypeMessageChannel{
			Data:    data,
			Message: "",
			//Param:   transmit.TPluginParam{},
		}
	})
	return nil
}

func (s *NatsInput) Uninstall(param ...interface{}) error {
	return nil
}

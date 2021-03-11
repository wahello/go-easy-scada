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

func NewNatsOutput() transmit.ISource {
	return &NatsOutput{}
}

type NatsOutput struct {
	addr  string
	token string
	topic string

	handleWrite chan transmit.TypeMessageChannel
}

func (s *NatsOutput) GetReader() chan transmit.TypeMessageChannel {

	return nil
}

func (s *NatsOutput) GetWriter() chan transmit.TypeMessageChannel {

	return s.handleWrite
}

func (s *NatsOutput) IsInputer() bool {
	return false
}

func (s *NatsOutput) IsOutputer() bool {
	return true
}

func (n *NatsOutput) SetToken(token string) {
	n.token = token
}

func (s *NatsOutput) Start(param ...interface{}) {
	utilio.AdapterInfoLog(s.token, pluginName, "启动服务")
}

func (s *NatsOutput) Install(param string) error {

	tempParams := map[string]*string{
		"addr":  &s.addr,
		"topic": &s.topic,
	}

	if err := utilio.ValidateConnStr(&tempParams, param); err != nil {
		return errors.New(utilio.GetAdapterLog(s.token, pluginName, err.Error()))
	}

	chans, err := mqInstance.NewOutput(s.token, s.addr, s.topic)
	if err != nil {
		return errors.New(utilio.GetAdapterLog(s.token, pluginName, err.Error()))
	}
	s.handleWrite = *chans
	return nil
}

func (s *NatsOutput) Uninstall(param ...interface{}) error {
	return nil
}

package natsmq

import (
	"errors"
	"esd-router-preview/adapter/transmit"
	"esd-router-preview/utilio"

	nats "github.com/nats-io/nats.go"
)

// MQ通讯包数据格式.
//type TypePacket struct {
//	MessageType string      `json:"message" mapstructure:"message"`
//	Data        interface{} `json:"data" mapstructure:"data"`
//	UUID        string      `json:"uuid" mapstructure:"uuid"`
//	Receiver    string      `json:"receiver" mapstructure:"receiver"`
//}

var pluginName = "mqtt"

type NatsService struct {
	conn *nats.EncodedConn

	recvTopic string
	sendTopic string
	addr      string
	token     string
	// 接口属性
	handleWrite chan transmit.TypeMessageChannel
	handleRead  chan transmit.TypeMessageChannel
}

func NewNatsService() transmit.ISource {
	return &NatsService{}
}

func (n *NatsService) SetToken(token string) {
	n.token = token
}

func (n *NatsService) sub(cb func(interface{})) error {
	if _, err := n.conn.Subscribe(n.recvTopic, cb); err != nil {
		return err
	}
	return nil
}

func (s *NatsService) GetReader() chan transmit.TypeMessageChannel {
	return s.handleRead
}

func (s *NatsService) GetWriter() chan transmit.TypeMessageChannel {
	return s.handleWrite
}

func (s *NatsService) Install(param string) error {

	tempParams := map[string]*string{
		"addr": &s.addr,
		"send": &s.sendTopic,
		"recv": &s.recvTopic,
	}

	if err := utilio.ValidateConnStr(&tempParams, param); err != nil {
		return errors.New(utilio.GetAdapterLog(s.token, pluginName, err.Error()))
	}

	s.handleWrite = make(chan transmit.TypeMessageChannel, 100)
	s.handleRead = make(chan transmit.TypeMessageChannel, 100)

	nc, err := nats.Connect(s.addr)
	if err != nil {
		return err
	}
	ec, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		return err
	}
	s.conn = ec
	if err = s.sub(func(data interface{}) {
		s.handleRead <- transmit.TypeMessageChannel{
			Data: data,
		}
	}); err != nil {
		return err
	}

	return nil
}

func (s *NatsService) Start(param ...interface{}) {
	// 等待发送函数
	go func() {
		utilio.AdapterInfoLog(s.token, pluginName, "启动服务")

		for {
			data := <-s.handleWrite
			err := s.conn.Publish(s.sendTopic, data.Data)
			if err != nil {
				utilio.AdapterErrLog(s.token, pluginName, err.Error())
			}
		}
	}()

}

func (s *NatsService) Uninstall(param ...interface{}) error {
	s.conn.Close()
	return nil
}

func (s *NatsService) IsInputer() bool {
	return true
}

func (s *NatsService) IsOutputer() bool {
	return true
}

func init() {
	transmit.Register("mqinput", NewNatsInput)
	transmit.Register("mqoutput", NewNatsOutput)
}

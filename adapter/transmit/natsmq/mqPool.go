/**
 * @Author: cmpeax Tang
 * @Date: 2021/3/12 1:53
 */

package natsmq

import (
	"esd-router-preview/adapter/transmit"
	"esd-router-preview/utilio"
	"github.com/nats-io/nats.go"
	"sync"
)

// 发送包结构
type TMQPACKET struct {
	Addr  string // 地址,token
	Topic string
	Data  interface{}
}

var mqInstance = &MQPool{
	pool: make(map[string]*MQClient, 0),
}

// 池化MQ
type MQPool struct {
	lock sync.RWMutex
	pool map[string]*MQClient
}

type MQClient struct {
	Conn        *nats.EncodedConn
	InputTopics []string

	InputCallbacks map[string]*[]func(interface{})
	OutputMap      map[string]*chan transmit.TypeMessageChannel
	Addr           string
	//token     string
	//// 接口属性
	//handleWrite chan MQClient

}

func (t *MQPool) NewConn(addr string) error {
	nc, err := nats.Connect(addr)
	if err != nil {
		return err
	}
	ec, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		return err
	}
	t.pool[addr] = &MQClient{
		Conn:        ec,
		Addr:        addr,
		InputTopics: []string{},

		InputCallbacks: make(map[string]*[]func(interface{}), 0),
		OutputMap:      make(map[string]*chan transmit.TypeMessageChannel, 0),
	}

	return nil
}

// 新创建 Input对象(不支持并发)
func (t *MQPool) NewInput(token string, addr string, topic string, callback func(interface{})) error {
	if _, exist := t.pool[addr]; !exist {
		// 新连接. 存入池并保存
		if err := t.NewConn(addr); err != nil {
			return err
		}
	}
	tconn := t.pool[addr]
	// 标记,寻找 topics 目录 检查是否已经有同样的topics
	flag := false
	for _, v := range tconn.InputTopics {
		if topic == v {
			flag = true
		}
	}

	// 寻找 map[topic] 回调数组是否已经创建,若没有,则创建
	if _, found := tconn.InputCallbacks[topic]; !found {
		tconn.InputCallbacks[topic] = &[]func(interface{}){callback}
	} else {
		*tconn.InputCallbacks[topic] = append(*tconn.InputCallbacks[topic], callback)
	}

	if flag == false {
		// 订阅
		if _, err := tconn.Conn.Subscribe(topic, func(data interface{}) {
			if _, ok := t.pool[addr]; !ok {
				utilio.AdapterErrLog(token, pluginName, "MQTT连接池里找不到该地址:"+addr)
				return
			}
			tmpconn := t.pool[addr]
			if _, ok := tmpconn.InputCallbacks[topic]; !ok {
				utilio.AdapterErrLog(token, pluginName, "MQTT连接池里找不到该订阅对象的回调."+topic)
				return
			}
			tconnCallback := tmpconn.InputCallbacks[topic]
			// 判断topic 若无,则警报.
			for _, tcb := range *tconnCallback {
				func(ttcb func(interface{})) {
					go ttcb(data)
				}(tcb)
			}
		}); err != nil {
			return err
		}
	}

	return nil
}

// 新创建 Input对象(不支持并发)
func (t *MQPool) NewOutput(token string, addr string, topic string) (*chan transmit.TypeMessageChannel, error) {
	if _, exist := t.pool[addr]; !exist {
		// 新连接. 存入池并保存
		if err := t.NewConn(addr); err != nil {
			return nil, err
		}
	}
	tconn := t.pool[addr]
	// 标记,寻找 topics 目录 检查是否已经有同样的topics

	if _, ok := tconn.OutputMap[topic]; !ok {
		tempPoint := make(chan transmit.TypeMessageChannel, 100)
		tconn.OutputMap[topic] = &tempPoint
		go func(ttopic string, chans chan transmit.TypeMessageChannel) {
			for {
				willSend := <-chans
				if err := t.pool[addr].Conn.Publish(topic, willSend.Data); err != nil {
					utilio.AdapterErrLog(token, pluginName, err.Error())
				}
			}
		}(topic, tempPoint)
		tconn.OutputMap[topic] = &tempPoint
	}

	return tconn.OutputMap[topic], nil
}

func (t *MQPool) Write(packet TMQPACKET) {

}

func (t *MQPool) Read(packet TMQPACKET) {

}

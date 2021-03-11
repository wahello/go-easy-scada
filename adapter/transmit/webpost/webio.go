package natsmq

import (
	"bytes"
	"errors"
	"esd-router-preview/adapter/transmit"
	"esd-router-preview/utilio"
	"net/http"
)

// MQ通讯包数据格式.
//type TypePacket struct {
//	MessageType string      `json:"message" mapstructure:"message"`
//	Data        interface{} `json:"data" mapstructure:"data"`
//	UUID        string      `json:"uuid" mapstructure:"uuid"`
//	Receiver    string      `json:"receiver" mapstructure:"receiver"`
//}
var pluginName = "webpost"

type WebPost struct {
	posturl string // 接收到信息后,post的url
	token   string
	// 接口属性
	handleWrite chan transmit.TypeMessageChannel
}

func NewWebPost() transmit.ISource {
	return &WebPost{}
}

func (n *WebPost) SetToken(token string) {
	n.token = token
}

func (s *WebPost) GetWriter() chan transmit.TypeMessageChannel {
	return s.handleWrite
}

func (s *WebPost) GetReader() chan transmit.TypeMessageChannel {
	return nil
}

func (s *WebPost) IsInputer() bool {
	return false
}

func (s *WebPost) IsOutputer() bool {
	return true
}

func (s *WebPost) Install(param string) error {

	tempParams := map[string]*string{
		"posturl": &s.posturl,
	}
	if err := utilio.ValidateConnStr(&tempParams, param); err != nil {
		return errors.New(utilio.GetAdapterLog(s.token, pluginName, err.Error()))
	}

	s.handleWrite = make(chan transmit.TypeMessageChannel, 100)

	return nil

}

func (s *WebPost) Start(param ...interface{}) {
	// 接收函数
	// 等待发送函数
	go func() {
		utilio.AdapterDebugLog(s.token, pluginName, "启动服务")

		for {
			data := <-s.handleWrite
			utilio.AdapterDebugLog(s.token, pluginName, "收到信息")
			res := data.Data.(string)

			// 异步请求，不阻塞
			go func(url string, tres []byte) {
				var tempClient = &http.Client{}
				var jsonStr = tres //转换二进制
				buffer := bytes.NewBuffer(jsonStr)

				request, err := http.NewRequest("POST", url, buffer)
				if err != nil {
					// 请求构造失败!
					utilio.AdapterErrLog(s.token, pluginName, err.Error())
					return
				}
				_, err = tempClient.Do(request)
				if err != nil {
					utilio.AdapterErrLog(s.token, pluginName, err.Error())
					return
				}

				// 后续可以根据条件 进行重试工作.

				// 如需要携带头部，使用该方法
				// request.Header.Set()

			}(s.posturl, []byte(res))

		}
	}()

}

func (s *WebPost) Uninstall(param ...interface{}) error {

	return nil
}

func init() {
	transmit.Register("webpost", NewWebPost)
}

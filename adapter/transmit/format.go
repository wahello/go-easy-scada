/**
 * @Author: cmpeax Tang
 * @Date: 2021/3/11 22:53
 */

package transmit

// 信息

type TPluginParam struct {
	PluginName string
	Data       interface{}
}

type TypeMessageChannel struct {
	Data    interface{}
	Message string
	Param   TPluginParam // 携带信息. 通过转换可检测
}

type ISource interface {
	IBase
	SetToken(string)
	IsInputer() bool
	IsOutputer() bool
	GetReader() chan TypeMessageChannel
	GetWriter() chan TypeMessageChannel
}

// 发布者数据源
type IInputSource interface {
	GetReader() chan TypeMessageChannel
}

// 订阅者写入源
type IOutputSource interface {
	GetWriter() chan TypeMessageChannel
}

type IBase interface {
	Install(string) error
	Start(...interface{})
	Uninstall(...interface{}) error
}

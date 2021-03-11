package logout

import (
	"esd-router-preview/core/logout/logg"
)

var L ILog = logg.NewLogg()

// 日志组件
type ILog interface {
	Set(string)
	ILogMethods
}

type ILogMethods interface {
	PrintLog(interface{}, ...interface{})
	Debug() interface{}
	Info() interface{}
	Error() interface{}
	Warn() interface{}
}

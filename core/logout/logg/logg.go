package logg

import (
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"os"
)

type LoggSys struct {
	logs *logrus.Logger
}

func NewLogg() *LoggSys {
	tempLog := &LoggSys{}
	tempLog.Init()
	return tempLog
}

const (
	// PanicLevel level, highest level of severity. Logs and then calls panic with the
	// message passed to Debug, Info, ...
	PanicLevel logrus.Level = iota
	// FatalLevel level. Logs and then calls `logger.Exit(1)`. It will exit even if the
	// logging level is set to Panic.
	FatalLevel
	// ErrorLevel level. Logs. Used for errors that should definitely be noted.
	// Commonly used for hooks to send errors to an error tracking service.
	ErrorLevel
	// WarnLevel level. Non-critical entries that deserve eyes.
	WarnLevel
	// InfoLevel level. General operational entries about what's going on inside the
	// application.
	InfoLevel
	// DebugLevel level. Usually only enabled when debugging. Very verbose logging.
	DebugLevel
	// TraceLevel level. Designates finer-grained informational events than the Debug.
	TraceLevel
)

func (t *LoggSys) Info() interface{} {
	return InfoLevel
}

func (t *LoggSys) Error() interface{} {
	return ErrorLevel
}

func (t *LoggSys) Warn() interface{} {
	return WarnLevel
}

func (t *LoggSys) Debug() interface{} {
	return DebugLevel
}

func (t *LoggSys) PrintLog(level interface{}, logarr ...interface{}) {
	t.logs.Logln(level.(logrus.Level), logarr)
	// fmt.Println(log)
	// logg.PrintLog(logg.DebugLevel,"[", strconv.Itoa(level), "]")
}

func (t *LoggSys) Init() {
	pathMap := lfshook.PathMap{}
	t.logs = logrus.New()
	t.logs.Hooks.Add(lfshook.NewHook(
		pathMap,
		&logrus.JSONFormatter{},
	))
	t.logs.SetFormatter(&logrus.JSONFormatter{})
	// 设置将日志输出到标准输出（默认的输出为stderr，标准错误）
	// 日志消息输出可以是任意的io.writer类型
	t.logs.SetOutput(os.Stdout)

	// 设置日志级别为warn以上
	tempLevel := InfoLevel
	//if levels, ok := param[0].(logrus.Level); ok {
	//	tempLevel = levels
	//} else {
	//	fmt.Println("[警告] 日志识别参数类型失败.默认选用Info输出.")
	//}

	t.logs.SetLevel(tempLevel)
}

func (t *LoggSys) Set(param string) {

}

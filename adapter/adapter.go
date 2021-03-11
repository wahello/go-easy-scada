package adapter

// 定义了接口的生命周期
type IBase interface {
	Install(string) error
	Start(...interface{})
	Uninstall(...interface{}) error
}

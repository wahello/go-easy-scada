/**
 * @Author: cmpeax Tang
 * @Date: 2021/3/10 16:23
 */

package adapter

// 类实现接口模板

type TBaseTemplate struct {
}

func (t *TBaseTemplate) Install(param string) error {
	return nil
}

func (t *TBaseTemplate) Start(param ...interface{}) error {
	return nil
}

func (t *TBaseTemplate) Uninstall(param ...interface{}) error {
	return nil
}

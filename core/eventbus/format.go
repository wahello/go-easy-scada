/**
 * @Author: cmpeax Tang
 * @Date: 2021/3/11 22:42
 */

package eventbus

import (
	"encoding/json"
	"errors"
)

type TypeDefineHook struct {
	Plugin  map[string]TypePlugin  `json:"plugin"`
	Message map[string]TypeMessage `json:"message"`
}

type TypePlugin struct {
	Module string `json:"module"`
	Param  string `json:"param"`
}

type TypeMessage struct {
	Pub []string `json:"pub"`
	Sub []string `json:"sub"`
}

// 底层组件
func NewEventBus(config string) (*TypeEventBusService, error) {
	// 转换 成json
	var tempHook TypeDefineHook
	if err := json.Unmarshal([]byte(config), &tempHook); err != nil {
		return nil, errors.New("Hook 描述文件初始化失败! " + err.Error())
	}
	return &TypeEventBusService{
		config: tempHook,
	}, nil
}

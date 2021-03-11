/**
 * @Author: cmpeax Tang
 * @Date: 2021/3/11 17:05
 */

package cache

import (
	"esd-router-preview/adapter"
	"fmt"
)

// 缓存组件
type ICache interface {
	adapter.IBase
	Get(interface{}) (interface{}, error)
	Set(interface{}, interface{}) error
}

type Instance func() ICache

var adapters = make(map[string]Instance)

func Register(name string, adapter Instance) {
	if adapter == nil {
		panic("cache: Register adapter is nil")
	}
	if _, ok := adapters[name]; ok {
		panic("cache: Register called twice for adapter " + name)
	}
	adapters[name] = adapter
}

func New(adapterName string, config string) (adapter ICache, err error) {
	instanceFunc, ok := adapters[adapterName]
	if !ok {
		err = fmt.Errorf("cache: unknown adapter name %q (forgot to import?)", adapterName)
		return nil, err
	}
	adapter = instanceFunc()
	return adapter, nil
}

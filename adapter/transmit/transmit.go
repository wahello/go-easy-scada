/**
 * @Author: cmpeax Tang
 * @Date: 2021/3/11 17:19
 */

package transmit

import (
	"fmt"
)

type Instance func() ISource

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

func New(adapterName string, config string, token string) (adapter ISource, err error) {
	instanceFunc, ok := adapters[adapterName]
	if !ok {
		err = fmt.Errorf("cache: unknown adapter name %q (forgot to import?)", adapterName)
		return nil, err
	}
	adapter = instanceFunc()
	adapter.SetToken(token)
	return adapter, adapter.Install(config)

}

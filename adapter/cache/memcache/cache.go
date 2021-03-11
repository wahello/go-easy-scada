/**
 * @Author: cmpeax Tang
 * @Date: 2021/3/10 16:23
 */

package memcache

import "sync"

type BaseCache struct {
	hash sync.Map
}

func NewBaseCache() *BaseCache {
	return &BaseCache{}
}

func (d *BaseCache) Set(k interface{}, v interface{}) {
	d.hash.Store(k, v)
}

func (d *BaseCache) Get(key interface{}) (interface{}, bool) {
	return d.hash.Load(key)
}

func (d *BaseCache) GetAll() map[interface{}]interface{} {
	result := make(map[interface{}]interface{}, 0)
	d.hash.Range(func(k interface{}, v interface{}) bool {
		result[k] = v
		return true
	})
	return result
}

func (t *BaseCache) Install(param ...interface{}) error {
	return nil
}

func (t *BaseCache) Start(param ...interface{}) error {
	return nil
}

func (t *BaseCache) Uninstall(param ...interface{}) error {
	return nil
}

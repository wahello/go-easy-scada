/**
 * @Author: cmpeax Tang
 * @Date: 2021/3/10 16:23
 */

package lockcache

import (
	"sync"
)

type LockData struct {
	lock sync.RWMutex
	data []interface{}
}

type LockCache struct {
	hash sync.Map
}

func NewLockCache() *LockCache {
	return &LockCache{}
}

// 修改
func (d *LockCache) Set(k interface{}, v func(interface{})) {
	value, _ := d.hash.LoadOrStore(k, &LockData{})
	value.(*LockData).lock.Lock()
	defer value.(*LockData).lock.Unlock()
	v(value.(*LockData).data)
}

func (d *LockCache) Get(key interface{}) (interface{}, bool) {
	realData, ok := d.hash.Load(key)

	if !ok {
		return nil, false
	}
	realData.(*LockData).lock.RLock()
	defer realData.(*LockData).lock.RUnlock()
	return realData.(*LockData).data, true

}

func (d *LockCache) Install(param ...interface{}) error {
	return nil
}

func (d *LockCache) Start(param ...interface{}) error {
	return nil
}

func (t *LockCache) Uninstall(param ...interface{}) error {
	return nil
}

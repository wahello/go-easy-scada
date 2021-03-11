/**
 * @Author: cmpeax Tang
 * @Date: 2021/3/11 11:39
 */

package core

// _ 会执行初始函数，避免 import cycle
// 可以通过 init 进行组件的选型.
import (
	"sync"
)

type GlobalComponent struct {
}

var instance *GlobalComponent
var once sync.Once

// 导入 描述通道型 JSON文件.
func GetInstance() *GlobalComponent {
	once.Do(func() {
		// 日志系统初始化
		instance = &GlobalComponent{}
	})
	return instance
}

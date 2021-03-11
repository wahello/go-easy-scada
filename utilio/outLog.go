/**
 * @Author: cmpeax Tang
 * @Date: 2021/3/12 1:40
 */

package utilio

import (
	"esd-router-preview/core/logout"
	"fmt"
)

func GetAdapterLog(token string, plugin string, msg string) string {
	return fmt.Sprintf("插件[%s(%s)]:%s", token, plugin, msg)
}

func AdapterErrLog(token string, plugin string, msg string) {
	logout.L.PrintLog(logout.L.Error(), GetAdapterLog(token, plugin, msg))
}
func AdapterDebugLog(token string, plugin string, msg string) {
	logout.L.PrintLog(logout.L.Debug(), GetAdapterLog(token, plugin, msg))
}

func AdapterInfoLog(token string, plugin string, msg string) {
	logout.L.PrintLog(logout.L.Info(), GetAdapterLog(token, plugin, msg))
}

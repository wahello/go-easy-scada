/**
 * @Author: cmpeax Tang
 * @Date: 2021/3/12 0:38
 */

package ticker

import (
	"esd-router-preview/core/tickermap/ctickermap"
)

type ITickerMap interface {
	Quit()
	Start()
}

func NewTickerMap(handler map[int64]func(), isCoroutines bool) ITickerMap {
	return ctickermap.NewTickerMap(handler, isCoroutines)
}

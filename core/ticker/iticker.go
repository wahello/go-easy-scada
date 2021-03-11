/**
 * @Author: cmpeax Tang
 * @Date: 2021/3/12 0:38
 */

package ticker

import (
	"esd-router-preview/core/ticker/cticker"
	"time"
)

type ITicker interface {
	Quit()
	Start()
}

type ITickerOnce interface {
	Start()
}

func NewTicker(pellet time.Duration, isCoroutines bool, handler func()) ITicker {
	return cticker.NewTicker(pellet, isCoroutines, handler)
}

func NewTickerOnce(pellet time.Duration, isCoroutines bool, handler func()) ITickerOnce {
	return cticker.NewTickerOnce(pellet, handler)
}

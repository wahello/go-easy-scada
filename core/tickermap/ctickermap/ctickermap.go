package ctickermap

import (
	"time"
)

type Ticker struct {
	handler      map[int64]func()
	quit         chan int
	isCoroutines bool
}

func NewTickerMap(handler map[int64]func(), isCoroutines bool) *Ticker {
	return &Ticker{
		handler:      handler,
		quit:         make(chan int, 0),
		isCoroutines: isCoroutines,
	}
}

func (t *Ticker) Quit() {
	t.quit <- 1
}

func (t *Ticker) Start() {
	startTime := time.Now().Unix() // 秒数
	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-ticker.C:
			for k, v := range t.handler {
				if (time.Now().Unix()-startTime)%k == 0 {
					if t.isCoroutines {
						go v()
					} else {
						v()
					}

				}

			}
		case <-t.quit:
			ticker.Stop()
			return
		}
	}

}

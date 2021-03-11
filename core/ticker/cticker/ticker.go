package cticker

import (
	"time"
)

type Ticker struct {
	isCoroutines bool
	handler      func()
	quit         chan int
	pellet       time.Duration
}

type TickerOnce struct {
	handler func()
	pellet  time.Duration
}

func NewTickerOnce(pellet time.Duration, handler func()) *TickerOnce {
	return &TickerOnce{
		handler: handler,
		pellet:  pellet,
	}
}

func (t *TickerOnce) Start() {
	ticker := time.NewTicker(t.pellet)
	<-ticker.C
	t.handler()

}

func NewTicker(pellet time.Duration, isCoroutines bool, handler func()) *Ticker {
	return &Ticker{
		handler:      handler,
		quit:         make(chan int, 0),
		pellet:       pellet,
		isCoroutines: isCoroutines,
	}
}

func (t *Ticker) Quit() {
	t.quit <- 1
}

func (t *Ticker) Start() {

	ticker := time.NewTicker(t.pellet)
	for {
		select {
		case <-ticker.C:
			if t.isCoroutines {
				go t.handler()
			} else {
				t.handler()
			}
		case <-t.quit:
			ticker.Stop()
			return
		}
	}

}

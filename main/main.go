package main

import (
	"esd-router-preview/engine"
	"fmt"
)

func main() {
	tengine := engine.GetInstance()
	tengine.LoadConfig("./interface.json")
	err := tengine.Init()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	//tengine.GetEventBus().Mount(custom.DealData)
	tengine.Start()
	//ticker.NewTicker(time.Second, true, func() {
	//
	//	tengine.GetEventBus().Send(transmit.TypeMessageChannel{
	//		Data:    custom.SendMonitor(),
	//		Message: "MESSAGE_SEND_TO_MONITOR",
	//	})
	//}).Start()

	for {

	}
	fmt.Println("部署完成")
}

package main

import (
	"github.com/jannchie/probe/centre"
	"github.com/jannchie/probe/client/collector"
)

func main() {
	engine := centre.InitCentre()
	addr := "localhost:12000"
	go collector.RunBaseCollector()
	_ = engine.Run(addr)
}

package main

import (
	"github.com/jannchie/probe/client/collector"
	"github.com/jannchie/probe/cmd"
)

func main() {
	engine := cmd.InitCentre()
	addr := "localhost:12000"
	go collector.RunBaseCollector()
	_ = engine.Run(addr)
}

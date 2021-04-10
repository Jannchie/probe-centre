/*
Copyright © 2020 Jannchie <jannchie@gmail.com>
*/
package main

import (
	"log"
	"os"

	"github.com/Jannchie/probe-centre/db"
	"github.com/Jannchie/probe-centre/job"
	"github.com/Jannchie/probe-centre/router"
)

func main() {
	log.SetFlags(log.Llongfile | log.LstdFlags)
	db.Init(os.Getenv("PROBE_PG_DSN"))
	go job.Init()
	_ = router.Init().Run(":12000")
}

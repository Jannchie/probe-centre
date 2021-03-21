/*
Copyright © 2020 Jannchie <jannchie@gmail.com>
*/
package main

import (
	"os"

	"github.com/Jannchie/probe-centre/db"
	"github.com/Jannchie/probe-centre/router"
)

func main() {
	db.Init(os.Getenv("PROBE_PG_DSN"))
	router.Init().Run(":12000")
}

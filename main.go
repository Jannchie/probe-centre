/*
Copyright © 2020 Jannchie <jannchie@gmail.com>
*/
package main

import (
	"github.com/Jannchie/pyobe-carrier/db"
	"github.com/Jannchie/pyobe-carrier/router"
)

func main() {
	db.InitDB()
	router.InitRouter().Run()
}

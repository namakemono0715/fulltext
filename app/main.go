package main

import (
	"fulltext/router"
)

func main() {
	r := router.SetupRouter()
	r.Run(":8080")
}
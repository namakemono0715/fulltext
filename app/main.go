package main

import (
	"fmt"
	"log"

	"fulltext/router"
	"github.com/joho/godotenv"
)

var index bleve.Index

func main() {
	// .envファイルの読み込み
	_ = godotenv.Load()
	r := router.SetupRouter()
	port := 8080
	fmt.Printf("Starting server on :%d...\n", port)
	log.Fatal(r.Run(fmt.Sprintf(":%d", port)))
}
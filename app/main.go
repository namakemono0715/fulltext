package main

import (
	"fmt"
	"log"
	"os"
	"runtime/debug"

	"fulltext/router"
	"github.com/joho/godotenv"
)

// var index bleve.Index

func main() {
	log.SetOutput(os.Stdout)
	log.Println("=== サーバーを起動します ===")
	
	defer func() {
    if r := recover(); r != nil {
      fmt.Println("Recovered from panic:", r)
      debug.PrintStack()
    }
  }()
	
	// .envファイルの読み込み
	_ = godotenv.Load()
	r := router.SetupRouter()
	port := 8080
	fmt.Printf("Starting server on :%d...\n", port)
	log.Fatal(r.Run(fmt.Sprintf(":%d", port)))
}
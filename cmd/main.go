package main

import (
	"log"
	"os"
)

// @title Movies Downloader API
// @version 1.0
// @description HTTP API for loading movie ranges and finding the last available movie.
func main() {
	err := run()
	if err != nil {
		log.Printf("server stopped with error: %v", err)
		os.Exit(1)
	}
}

package main

import (
	"log"
	"os"

	"claude-proxy/zedmode"
	"claude-proxy/httpmode"
)

func main() {
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		zedmode.Run(os.Stdin, os.Stdout, os.Stderr)
	} else {
		log.Println("Starting HTTP server on :8080")
		httpmode.Start()
	}
}

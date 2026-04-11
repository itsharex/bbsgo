//go:build !gui

package main

import (
	"fmt"
	"log"
)

func main() {
	ip, port := ParseFlags()
	Init()

	addr := fmt.Sprintf("%s:%d", ip, port)
	srv := NewServer(addr)

	if err := srv.Start(); err != nil {
		log.Fatalf("server failed to start: %v", err)
	}
}

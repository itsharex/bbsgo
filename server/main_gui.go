//go:build gui

package main

import (
	"bbsgo/systray"
	"context"
	"fmt"
	"log"
)

func main() {
	ip, port := ParseFlags()
	Init()

	addr := fmt.Sprintf("%s:%d", ip, port)
	srv := NewServer(addr)

	shutdownFunc := func(ctx context.Context) error {
		return srv.Shutdown(ctx)
	}

	getActualAddrFunc := func() string {
		return addr
	}

	systray.Init(addr, shutdownFunc, getActualAddrFunc)

	systray.Run(func() {
		systray.Setup()
		go func() {
			if err := srv.Start(); err != nil {
				log.Printf("server error: %v", err)
			}
		}()
	})
}

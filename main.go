package main

import (
	"KamaiZen/logger"
	"KamaiZen/lsp"
	"KamaiZen/server"
	"KamaiZen/settings"
	"KamaiZen/state_manager"
	"flag"
	"fmt"
	"sync"
)

func initialise() {
	logger.Info("Starting KamaiZen")
	state_manager.InitializeState()
	lsp.Initialise()
}

func main() {
	v := flag.Bool("version", false, "print version")
	flag.Parse()
	if *v {
		fmt.Printf("version %s\n", settings.KAMAIZEN_VERSION)
		return
	}
	initialise()
	defer logger.Info("KamaiZen stopped")
	server := server.GetServerInstance()

	var wg sync.WaitGroup
	wg.Add(2)
	go server.StartServer(&wg)
	go lsp.Start(&wg)
	wg.Wait()
}

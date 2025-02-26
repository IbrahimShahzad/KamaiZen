package main

import (
	"KamaiZen/lsp"
	"KamaiZen/server"
	"KamaiZen/settings"
	"KamaiZen/state_manager"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	v := flag.Bool("version", false, "print version")
	flag.Parse()
	if *v {
		fmt.Printf("version %s\n", settings.KAMAIZEN_VERSION)
		return
	}
	initialise()
	defer log.Info().Msg("KamaiZen stopped")
	server := server.GetServerInstance()

	var wg sync.WaitGroup
	wg.Add(2)
	go server.StartServer(&wg)
	go lsp.Start(&wg)
	wg.Wait()
}

func initialise() {
	lev := zerolog.InfoLevel
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMicro
	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		return filepath.Base(file) + ":" + strconv.Itoa(line)
	}
	file, err := os.OpenFile(
		"/tmp/kamaizen.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0664,
	)
	if err != nil {
		panic(err)
	}
	zerolog.TimeFieldFormat = zerolog.TimestampFunc().UTC().Format("2006-01-02T15:04:05.000Z")
	log.Logger = zerolog.New(file).With().Caller().Timestamp().Logger().Level(lev)
	log.Info().Msg("Starting KamaiZen!")
	state_manager.InitializeState()
	lsp.Initialise()
}

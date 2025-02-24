package lsp

import (
	"KamaiZen/rpc"
	"github.com/rs/zerolog/log"
	"io"
	"os"
	"sync"
)

var writer io.Writer

const buffered_channel_size = 24

var writer_channel = make(chan []byte, buffered_channel_size)

// WriteResponse encodes the given response and sends it to the writer channel.
//
// Parameters:
//
//	response interface{} - The response to be encoded and written.
func WriteResponse(response interface{}) {
	reply := rpc.EncodeMessage(response)
	writer_channel <- []byte(reply)
}

// Write writes the given message to the writer.
//
// Parameters:
//
//	message []byte - The message to be written.
func Write(message []byte) {
	writer.Write(message)
}

// Start starts the writer goroutine that listens for messages on the writer channel
// and writes them to the writer. It signals the wait group when done.
//
// Parameters:
//
//	wg *sync.WaitGroup - The wait group to signal when done.
func Start(wg *sync.WaitGroup) {
	defer wg.Done()
	log.Info().Msg("Starting writer")
	for {
		select {
		case message := <-writer_channel:
			Write(message)
		}
	}
}

// Initialise initializes the writer to use os.Stdout.
func Initialise() {
	writer = os.Stdout
}

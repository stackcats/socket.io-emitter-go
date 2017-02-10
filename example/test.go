package main

import (
	"github.com/stackcats/socket.io-emitter-go"
)

func main() {
	opts := &emitter.Options{}
	socket := emitter.NewEmitter(opts)
	defer socket.Close()
	socket.Broadcast().Emit("ping", "Hello World")
}

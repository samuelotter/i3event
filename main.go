package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/samuelotter/i3ipc"
)

var configFile = flag.String("config", "~/.i3event", "Path to config file.")
var debug      = flag.Bool("debug", false, "Activate debug logging.")

func Debugf(format string, args ...interface{}) {
	if *debug {
		log.Printf("DEBUG: %s", fmt.Sprintf(format, args...))
	}
}

func main() {
	flag.Parse()

	config, err := ReadConfiguration(*configFile)
	if err != nil {
		log.Fatalf("Failed to open configuration file: %v\n", err)
	}

	channel := make(chan i3ipc.Event)
	SubscribeChannel(i3ipc.I3WorkspaceEvent, channel)
	SubscribeChannel(i3ipc.I3OutputEvent, channel)
	SubscribeChannel(i3ipc.I3ModeEvent, channel)
	SubscribeChannel(i3ipc.I3WindowEvent, channel)
	SubscribeChannel(i3ipc.I3BarConfigEvent, channel)

	EventLoop(channel, config)
}

func SubscribeChannel(eventType i3ipc.EventType, aggregate chan i3ipc.Event) (err error) {
	channel, err := i3ipc.Subscribe(eventType)
	if err != nil {
		return
	}
	go func() {
		for event := range channel {
			aggregate <- event
		}
	}()
	return
}

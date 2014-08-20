package main

import (
	"encoding/json"
	"fmt"

	"os"
	"os/exec"

	"github.com/samuelotter/i3ipc"
)

var eventTypes = map[string]i3ipc.EventType {
	"workspace": i3ipc.I3WindowEvent,
	"output": i3ipc.I3OutputEvent,
	"mode": i3ipc.I3ModeEvent,
	"window": i3ipc.I3WindowEvent,
	"barupdate_config": i3ipc.I3BarConfigEvent,
}

func EventLoop(events chan i3ipc.Event, config *Config) {
	eventMap := make(map[i3ipc.EventType][]Rule)
	for _, rule := range config.Rules {
		eventMap[rule.Event] = append(eventMap[rule.Event], rule)
	}

        for event := range events {
		Debugf("Received event: %+v", event)
		rules := eventMap[event.Type]
		for _, rule := range rules {
			if event.Change == rule.Change || rule.Change == "*" {
				Debugf("Found matching rule %+v\n", rule)
				switch rule.Action {
				case ActionExec:
					cmd := exec.Command(rule.Args[0],
						rule.Args[1:]...)
					cmd.Stdout = os.Stdout
					cmd.Stderr = os.Stderr
					StdIn, err := cmd.StdinPipe()
					if err != nil {
						fmt.Println(err)
						return
					}
					if err := cmd.Start(); err != nil {
						fmt.Println(err)
						return
					}
					msg, err := json.Marshal(event.Payload)
					if err != nil {
						fmt.Println(err)
						return
					}
					Debugf("STDIN << %s\n", msg)
					StdIn.Write(msg)
					StdIn.Close()
					cmd.Wait()
				default:
					// Ignore
				}
			}
		}
	}
}

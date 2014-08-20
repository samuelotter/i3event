package main

import (
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
			if rule.Match(event) {
				Debugf("Found matching rule %+v\n", rule)
				rule.Handle(event)
			}
		}
	}
}

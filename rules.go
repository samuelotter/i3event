package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"github.com/samuelotter/i3ipc"
)

type Rule struct {
	Event   i3ipc.EventType
	Change  string
	Action  Action
}

func (self *Rule) Match(event i3ipc.Event) bool {
	return event.Change == self.Change || self.Change == "*"
}

func (self *Rule) Handle(event i3ipc.Event) {
	self.Action.Invoke(event)
}

type Action interface {
	Invoke(event i3ipc.Event) error
}

type IgnoreAction struct {

}

func (self IgnoreAction) Invoke(event i3ipc.Event) error {
	Debugf("Ignored event %+v\n", event)
	return nil
}

type ExecAction struct {
	Args []string
}

func (self *ExecAction) Invoke(event i3ipc.Event) error {
	Debugf("ExecAction.Invoke %+v", self)
	cmd := exec.Command(self.Args[0], self.Args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	StdIn, err := cmd.StdinPipe()
	if err != nil {
		fmt.Println(err)
		return err
	}
	if err := cmd.Start(); err != nil {
		fmt.Println(err)
		return err
	}
	msg, err := json.Marshal(event.Payload)
	if err != nil {
		fmt.Println(err)
		return err
	}
	Debugf("STDIN << %s\n", msg)
	StdIn.Write(msg)
	StdIn.Close()
	cmd.Wait()
	return nil
}

func NewAction(name string, args []string) Action {
	switch name {
	case "ignore":
		return &IgnoreAction{}
	case "exec":
		return &ExecAction{Args: args}
	}
	log.Panic("Unsupported action")
	return nil
}

package handler

import (
	"bytes"
	"os/exec"
	"strings"

	"github.com/hevnly/eevy/event"
	"github.com/hevnly/eevy/handler/config"
)

// This listener performs system calls. Perhaps a local binary needs to be called, use this listener.
type Cli struct {
	HandlerBase

	Config config.Cli
}

// Satisfies the Listener interface and calls the relavent binary file
func (this *Cli) Exec(evt event.Event) {

	this.Log.Handler(this, &evt)

	bin := magicString(this.Config.GetBin(), evt)
	stdin := magicString(this.Config.GetStdin(), evt)
	args := append([]string(nil), this.Config.GetArgs()...)
	for i, arg := range this.Config.GetArgs() {
		args[i] = magicString(arg, evt)
	}

	cmd := exec.Command(bin, args...)
	cmd.Stdin = strings.NewReader(stdin)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err := cmd.Run()
	if err != nil {
		this.Log.Error("%s: %s", out.String(), err.Error())
		return
	}
}

func (this *Cli) GetType() string {

	return this.GetConfig().GetType()
}

func (this *Cli) GetConfig() config.Handler {
	return this.Config
}

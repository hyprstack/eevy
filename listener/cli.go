package listener

import (
	"bytes"
	"os/exec"
	"strings"

	"github.com/hevnly/eevy/event"
)

type CliConfig interface {
	ListenerConfig

	GetBin() string
	GetArgs() []string
	GetStdin() string
}

// This listener performs system calls. Perhaps a local binary needs to be called, use this listener.
type Cli struct {
	ListenerBase

	Config CliConfig
}

// Satisfies the Listener interface and calls the relavent binary file
func (this *Cli) Exec(evt event.Event) {

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
		gLog.Error("%s: %s", out.String(), err.Error())
		return
	}
	gLog.Debug("%s", out.String())
}

package listener

import (
	"bytes"
	"os/exec"
	"strings"

	"github.com/hevnly/eevy/config"
	"github.com/hevnly/eevy/event"
)

type Cli struct {
	ListenerBase

	Bin  string
	Args []string
}

func (this *Cli) Exec(evt event.Event) {

	bin := this.magicString(this.Bin, evt)
	args := this.Args

	for i, arg := range this.Args {
		args[i] = this.magicString(arg, evt)
	}
	gLog.Info("here: %s %q", bin, args)
	cmd := exec.Command(bin, args...)
	cmd.Stdin = strings.NewReader("some input")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	if err != nil {
		gLog.Error("%s: %s", out.String(), err.Error())
		return
	}
	gLog.Debug("%q", out.String())
}

func (this *Cli) Init(conf config.Listener) {
	this.Bin = conf.Bin
	this.Args = conf.Args
}

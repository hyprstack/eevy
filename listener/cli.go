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

	Bin   string
	Args  []string
	Stdin string
}

func (this *Cli) Exec(evt event.Event) {

	bin := this.magicString(this.Bin, evt)
	stdin := this.magicString(this.Stdin, evt)
	args := append([]string(nil), this.Args...)
	for i, arg := range this.Args {
		args[i] = this.magicString(arg, evt)
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

func (this *Cli) Init(conf config.Listener) {
	this.Bin = conf.Bin
	this.Args = conf.Args
	this.Stdin = conf.Stdin
}

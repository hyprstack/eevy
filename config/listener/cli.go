package listener

type Cli struct {
	ListenerBase
}

func (this *Cli) GetArgs() []string {
	return this.getStringSlice("args")
}

func (this *Cli) GetBin() string {
	return this.getString("bin")
}

func (this *Cli) GetStdin() string {
	return this.getString("stdin")
}

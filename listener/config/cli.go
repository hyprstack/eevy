package config

type Cli interface {
	Listener

	GetBin() string
	GetArgs() []string
	GetStdin() string
}

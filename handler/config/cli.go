package config

type Cli interface {
	Handler

	GetBin() string
	GetArgs() []string
	GetStdin() string
}

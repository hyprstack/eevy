package config

type Lambda interface {
	Listener

	GetFunction() string
	GetRegion() string
}

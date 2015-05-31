package config

type Lambda interface {
	Handler

	GetFunction() string
	GetRegion() string
}

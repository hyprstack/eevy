package config

type Sqs interface {
	Listener

	GetUrl() string
	GetRegion() string
}

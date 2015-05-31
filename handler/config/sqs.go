package config

type Sqs interface {
	Handler

	GetUrl() string
	GetRegion() string
}

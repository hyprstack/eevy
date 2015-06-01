package config

type OAuth2 interface {
	Handler

	GetClientId() string
	GetClientSecret() string
	GetTokenUrl() string
	GetScope() []string
	GetEndPoint() string
	GetVerb() string
	GetBody() string
}

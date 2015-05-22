package config

type Listener struct {
	Type	     string
	Function     string
	Url          string
	Region       string
	ClientId     string    `yaml:"clientId"`
	ClientSecret string    `yaml:"clientSecret"`
	TokenUrl     string    `yaml:"tokenUrl"`
	Scope        []string
	EndPoint     string    `yaml:"endPoint"`
	Message		 string
	Verb		 string
}

package config

type Source struct {
	Type   string
	Host   string
	Port   int
	BindIp string `yaml:"bindIp"`
	Url    string
	Region string
}

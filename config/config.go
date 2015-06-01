// Various stucts that are used to configure sources and handlers.
package config

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// Represents the structure that should be used in configuration files
type Config struct {
	Sources   []Source
	Handlers  HandlerList
	Listeners ListenerList
}

type ListenerList map[string][]string

type HandlerList map[string]Handler

type Handler map[string]interface{}

func (this *Handler) get(s string) string {

	if _, ok := (*this)[s]; !ok {
		return ""
	}
	if t, ok := (*this)[s].(string); ok {
		return t
	}
	return ""
}

// Return the type of handler eg "sqs", "lambda" etc
func (this *Handler) GetType() string {
	return this.get("type")
}

// Return the message that should be sent by the handler, default is "${message}"
func (this *Handler) GetMessage() string {

	r := this.get("message")
	if r != "" {
		return r
	}
	return "${message}"
}

//  Convert this handler into a json string
func (this *Handler) String() string {
	b, _ := json.Marshal(this)
	return string(b)
}

func (this *Handler) Init(s string) {
	return
}

// Build this configuration struct from a configuration file
func (this *Config) LoadFromFile(s string) {

	filename, _ := filepath.Abs(s)
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(yamlFile, &this)
	if err != nil {
		panic(err)
	}
}

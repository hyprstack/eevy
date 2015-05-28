// Various stucts that are used to configure sources and listeners.
package config

import (
	"encoding/json"
)

type Config struct {
	Sources   []Source
	Listeners ListenerList
}

type ListenerList map[string][]Listener

type Listener map[string]interface{}

func (this *Listener) GetType() string {

	if _, ok := (*this)["type"]; !ok {
		return ""
	}
	if t, ok := (*this)["type"].(string); ok {
		return t
	}
	return ""
}

func (this *Listener) GetMessage() string {

	if _, ok := (*this)["message"]; !ok {
		return ""
	}
	if t, ok := (*this)["message"].(string); ok {
		return t
	}
	return "${message}"
}

func (this *Listener) String() string {
	b, _ := json.Marshal(this)
	return string(b)
}

func (this *Listener) Init(s string) {
	return
}

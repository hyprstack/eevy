package listener

import (
	"encoding/json"

	"gopkg.in/yaml.v2"
)

type ListenerBase struct {
	options map[string]interface{}
}

func (this *ListenerBase) Init(s string) {

	err := json.Unmarshal([]byte(s), &this.options)
	if err == nil {
		return
	}

	yaml.Unmarshal([]byte(s), &this.options)
	return
}

func (this *ListenerBase) String() string {
	s, _ := json.Marshal(this.options)
	return string(s)
}

func (this *ListenerBase) get(s string) interface{} {

	if b, ok := this.options[s]; ok {
		return b
	}
	return nil
}

func (this *ListenerBase) getString(s string) string {

	i := this.get(s)
	if m, ok := i.(string); ok {
		return m
	}
	return ""
}

// TODO(100% sure this could be done better)
func (this *ListenerBase) getStringSlice(s string) []string {
	a := this.get(s)
	if m, ok := a.([]interface{}); ok {
		sl := make([]string, len(m))
		for i, v := range m {
			if s, ok := v.(string); ok {
				sl[i] = string(s)
			}
		}
		return sl
	}
	return nil
}

func (this *ListenerBase) GetMessage() string {

	i := this.get("message")
	if i == nil {
		return "${message}"
	}
	if m, ok := i.(string); ok {
		return m
	}
	return ""
}

func (this *ListenerBase) GetType() string {

	return this.getString("type")
}

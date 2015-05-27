// Events are created by sources when the system they are listening on triger an event.
package event

import (
	"strings"

	"github.com/op/go-logging"
)

type Event struct {
	Event   string                 `json:"event"`
	Message map[string]interface{} `json:"message,omitempty"`
	Id      string                 `json:"id,omitempty"`
	Error   string                 `json:"error,omitempty"`
}

var gLog = logging.MustGetLogger("applog")

func (this *Event) Get(s string) interface{} {

	return get(s, this.Message)
}

func (this *Event) GetString(s string) string {

	val := this.Get(s)
	r, _ := val.(string)
	return r
}

func (this *Event) GetInt(s string) int {

	val := this.Get(s)
	r, ok := val.(float64)
	if ok {
		return int(r)
	}
	return 0
}

func (this *Event) GetFloat(s string) float64 {

	val := this.Get(s)
	return val.(float64)
}

func get(s string, val interface{}) interface{} {

	w, ok := val.(map[string]interface{})
	if !ok {
		return val
	}

	ns := strings.Split(s, ".")
	if sub, ok := w[ns[0]]; ok {
		return get(strings.Join(ns[1:], "."), sub)
	}
	return nil
}

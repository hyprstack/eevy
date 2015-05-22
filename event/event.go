package event

import (
	"encoding/json"
	"strings"
)

type Event struct {
	Event    string `json:"event"`
	Message  string `json:"message,omitempty"`
	Id       string `json:"id,omitempty"`
	Error    string `json:"error,omitempty"`

	messageJson MessageObj
}

type MessageObj map[string]interface{}

func (this *Event) MessageJson() MessageObj {

	if this.messageJson != nil {
		return this.messageJson
	}

	rtn := make(MessageObj)
	json.Unmarshal([]byte(this.Message), &rtn)

	this.messageJson = rtn
	return rtn
}

func (this *Event) Get(s string) string {

	msg := this.MessageJson()
	return msg.get(s)
}

func (this *MessageObj) get(s string) string {

	opt := strings.Split(s, ".")

	if sub, ok := (*this)[opt[0]]; ok {

		if str, ok := sub.(string); ok {
			return str
		}
	}

	return ""
}

// Listeners that can be used when events are picked up by any of the sources.
package listener

import (
	"regexp"
	"strings"
	"sync"

	"encoding/json"
	"github.com/op/go-logging"

	"github.com/hevnly/eevy/config"
	listConfig "github.com/hevnly/eevy/config/listener"
	"github.com/hevnly/eevy/event"
)

// Think of this as a collection of listeners grouped by the Event name
// This also stores all of the child EventListners creating a tree structure
type EventListener struct {
	Name      string
	Listeners []Listener
	Subset    map[string]*EventListener

	wildCard bool
	subLock  sync.Mutex
}

type Listener interface {
	Exec(evt event.Event)
}

type ListenerBase struct {
}

type ListenerConfig interface {
	GetType() string
	GetMessage() string
	String() string
	Init(s string)
}

var gLog = logging.MustGetLogger("applog")

// Recieves a configuration struct and creates the relavent Listener
func BuildFromConf(conf ListenerConfig) Listener {

	var l Listener
	switch conf.GetType() {
	case "sqs":
		tl := new(Sqs)
		tc := new(listConfig.Sqs)
		tc.Init(conf.String())
		tl.Config = tc
		l = tl

	case "lambda":
		tl := new(Lambda)
		tc := new(listConfig.Lambda)
		tc.Init(conf.String())
		tl.Config = tc
		l = tl

	case "oauth2":
		tl := new(OAuth2)
		tc := new(listConfig.OAuth2)
		tc.Init(conf.String())
		tl.Config = tc
		l = tl

	case "cli":
		tl := new(Cli)
		tc := new(listConfig.Cli)
		tc.Init(conf.String())
		tl.Config = tc
		l = tl

	}
	return l
}

// Add the Listener to this EventListener using the supplied event name (evt eg "test.subTest")
func (l *EventListener) Add(evt string, lst Listener) {

	if evt == "_" {
		evt = "*"
	}

	// if this is the event we are trying to add to
	if evt == l.Name {
		l.Listeners = append(l.Listeners, lst)
		return
	}

	name := strings.Replace(evt, l.Name, "", -1)
	ns := strings.Split(name, ".")
	// if the first name space is empty remove it
	if ns[0] == "" {
		ns = append(ns[:0], ns[1:]...)
	}
	leng := len(ns)

	// if this is the event we are trying to add to
	// but we had a trailing '.' on the name
	if leng == 0 {
		l.Listeners = append(l.Listeners, lst)
		return
	}

	// if we have a subset for this event already
	if sub, ok := l.Subset[ns[0]]; ok {
		sub.Add(evt, lst)
		return
	}

	// we got this far so we need to create a new subset
	newName := ""
	if l.Name != "" {
		newName = l.Name + "." + ns[0]
	} else {
		newName = ns[0]
	}
	newSub := EventListener{
		Name:     newName,
		wildCard: ns[0] == "*",
	}
	newSub.Add(evt, lst)

	if l.Subset == nil {
		l.Subset = make(map[string]*EventListener)
	}
	l.Subset[ns[0]] = &newSub
}

// Executes the given event.
func (l *EventListener) Exec(evt event.Event) {

	// work out the event name relative to this listener
	relName := strings.Replace(evt.Event, l.Name, "", -1)

	l.subLock.Lock()
	// execute all listners that end here
	if l.wildCard == true || l.Name == evt.Event {
		for _, listener := range l.Listeners {
			listener.Exec(evt)
		}
	}
	// execute the wildcard listeners
	if sub, ok := l.Subset["*"]; ok {
		sub.Exec(evt)
	}

	// return if we have reached the end of the event name
	if relName == "" {
		l.subLock.Unlock()
		return
	}

	ns := strings.Split(relName, ".")
	if ns[0] == "" {
		ns = append(ns[:0], ns[1:]...)
	}

	// dive deep into the subset if one exists
	if sub, ok := l.Subset[ns[0]]; ok {
		sub.Exec(evt)
	}
	l.subLock.Unlock()
}

func BuildListeners(conf *config.ListenerList) *EventListener {
	rootListener := EventListener{}
	rootListener.Name = ""
	for evtName, listners := range *conf {
		for _, l := range listners {

			list := BuildFromConf(&l)
			rootListener.Add(evtName, list)
		}
	}
	return &rootListener
}

// Replaces variables ("${}") in the string to their actual value
func magicString(s string, evt event.Event) string {

	rep := regexp.MustCompile("(\\${|})")
	rst := findMagicStrings(s)
	for _, v := range rst {
		variable := rep.ReplaceAllString(v, "")
		opt := strings.Split(variable, ".")
		switch opt[0] {
		case "message":
			str := ""
			if len(opt) <= 1 {
				b, _ := json.Marshal(evt.Message)
				str = string(b)
			} else {
				str = evt.GetString(strings.Join(opt[1:], "."))
			}
			s = strings.Replace(s, v, str, -1)
		}
	}
	return s
}

// Finds all of the ${} in the given string
func findMagicStrings(s string) []string {
	re := regexp.MustCompile("\\${(.*?)}")
	return re.FindAllString(s, -1)
}

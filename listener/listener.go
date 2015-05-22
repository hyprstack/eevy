package listener

import (
	"regexp"
	"strings"
	"sync"

	"github.com/op/go-logging"

	"hevnly/eevy/config"
	"hevnly/eevy/event"
)

type EventListener struct {
	Name      string
	Listeners []Listener
	Subset    map[string]*EventListener

	wildCard bool
	subLock  sync.Mutex
}

type Listener interface {
	Exec(evt event.Event)
	Init(conf config.Listener)

	GetMessage(evt event.Event) string
}

type ListenerBase struct {
	Message string
}

var gLog = logging.MustGetLogger("applog")

func BuildFromConf(conf config.Listener) Listener {
	var list Listener
	switch conf.Type {
	case "sqs":
		list = new(Sqs)
	case "lambda":
		list = new(Lambda)
	case "oauth2":
		list = new(OAuth2)
	}
	list.Init(conf)
	return list
}

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

func (this *ListenerBase) magicString(s string, evt event.Event) string {

	rep := regexp.MustCompile("(\\${|})")
	rst := this.findMagicStrings(s)
	for _, v := range rst {
		variable := rep.ReplaceAllString(v, "")
		opt := strings.Split(variable, ".")
		switch opt[0] {
		case "message":
			str := ""
			if len(opt) <= 1 {
				str = evt.Message
			} else {
				str = evt.Get(strings.Join(opt[1:], "."))
			}
			s = strings.Replace(s, v, str, -1)
		}
	}
	return s
}

func (this *ListenerBase) findMagicStrings(s string) []string {
	re := regexp.MustCompile("\\${(.*?)}")
	return re.FindAllString(s, -1)
}

func (this *ListenerBase) GetMessage(evt event.Event) string {
	if this.Message == "" || this.Message == "${message}" {
		return evt.Message
	}
	return this.magicString(this.Message, evt)
}

// Listens on sources and envokes handlers when relavent events are detected
package listener

import (
	"strings"
	"sync"

	"github.com/hevnly/eevy/config"
	"github.com/hevnly/eevy/event"
	"github.com/hevnly/eevy/handler"
	"github.com/hevnly/eevy/logger"
)

// Think of this as a collection of listeners grouped by the Event name
// This also stores all of the child Handlers creating a tree structure
type Listener struct {
	Name     string
	Handlers []handler.Handler
	Subset   map[string]*Listener

	wildCard bool
	subLock  sync.Mutex
}

func BuildListener(conf config.ListenerList, hl *handler.HandlerList, log logger.Logger) *Listener {
	rootListener := Listener{}
	rootListener.Name = ""
	for evtName, listners := range conf {
		for _, l := range listners {

			h := hl.Get(l)
			rootListener.Add(evtName, h)
		}
	}
	return &rootListener
}

// Add the Handler to this Listener using the supplied event name (evt eg "test.subTest")
func (l *Listener) Add(evt string, handler handler.Handler) {

	if evt == "_" {
		evt = "*"
	}

	// if this is the event we are trying to add to
	if evt == l.Name {
		l.Handlers = append(l.Handlers, handler)
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
		l.Handlers = append(l.Handlers, handler)
		return
	}

	// if we have a subset for this event already
	if sub, ok := l.Subset[ns[0]]; ok {
		sub.Add(evt, handler)
		return
	}

	// we got this far so we need to create a new subset
	newName := ""
	if l.Name != "" {
		newName = l.Name + "." + ns[0]
	} else {
		newName = ns[0]
	}
	newSub := Listener{
		Name:     newName,
		wildCard: ns[0] == "*",
	}
	newSub.Add(evt, handler)

	if l.Subset == nil {
		l.Subset = make(map[string]*Listener)
	}
	l.Subset[ns[0]] = &newSub
}

// Executes the given event.
func (l *Listener) Exec(evt event.Event) {

	// work out the event name relative to this listener
	relName := strings.Replace(evt.Event, l.Name, "", -1)

	l.subLock.Lock()
	// execute all handlers that end here
	if l.wildCard == true || l.Name == evt.Event {
		for _, handler := range l.Handlers {
			handler.Exec(evt)
		}
	}
	// execute the wildcard handlers
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

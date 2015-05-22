package source

import (
	"crypto/rand"
	"encoding/base64"
	"sync"

	"encoding/json"

	"github.com/op/go-logging"

	"hevnly/eevy/config"
	"hevnly/eevy/event"
	"hevnly/eevy/listener"
)

func BuildFromConfig(conf config.Source, rootList *listener.EventListener) *Source {

	appLog := logging.MustGetLogger("applog")

	var src Source
	switch conf.Type {
	case "http":
		src = new(Http)
	case "sqs":
		src = new(Sqs)
	}
	src.init(appLog, conf, rootList)

	return &src
}

type Source interface {
	Listen(wg sync.WaitGroup)
	init(log *logging.Logger, conf config.Source, rootList *listener.EventListener)
}

type Base struct {
	config.Source

	Listener *listener.EventListener
	listLock sync.Mutex

	AppLog *logging.Logger
}

func (s *Base) init(log *logging.Logger, conf config.Source, rootList *listener.EventListener) {
	s.AppLog = log
	s.Url = conf.Url
	s.Region = conf.Region
	s.Port = conf.Port
	s.Listener = rootList
}

func (s *Base) processRaw(msg string) event.Event {

	var evt event.Event
	if err := json.Unmarshal([]byte(msg), &evt); err != nil {
		s.AppLog.Error("Can not turn string into event: %s", msg)
		return evt
	}
	evt.Id = generateId()
	s.AppLog.Info("Event created. Id=%s event=%s", evt.Id, evt.Event)
	s.processEvent(evt)

	return evt
}

func (s *Base) processEvent(evt event.Event) {

	s.listLock.Lock()
	s.Listener.Exec(evt)
	s.listLock.Unlock()
}

func generateId() string {

	size := 32 // change the length of the generated random string here
	rb := make([]byte, size)
	_, err := rand.Read(rb)
	if err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(rb)
}

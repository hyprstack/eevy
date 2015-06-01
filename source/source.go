// Sources are are where events are placed or called from.  They can listen
// on a specific http port, watch an AWS SQS queue etc
package source

import (
	"crypto/rand"
	"encoding/base64"
	"sync"

	"encoding/json"

	"github.com/hevnly/eevy/config"
	"github.com/hevnly/eevy/event"
	"github.com/hevnly/eevy/listener"
	"github.com/hevnly/eevy/logger"
)

// Recieve a configuration struct and create the relevant source
func BuildFromConfig(conf config.Source, rootList *listener.Listener, log logger.Logger) Source {

	var src Source
	switch conf.Type {
	case "http":
		src = &Http{}
	case "sqs":
		src = &Sqs{}
	}
	src.init(log, conf, rootList)

	return src
}

// Pass in the source configuration and this function both builds and starts listening to the source
func StartSources(sourceConf *[]config.Source, rootList *listener.Listener, log logger.Logger, wg sync.WaitGroup) {

	var wgLocal sync.WaitGroup
	var sources []Source
	for _, conf := range *sourceConf {
		tmp := BuildFromConfig(conf, rootList, log)

		sources = append(sources, tmp)
		wgLocal.Add(1)
		go tmp.Listen(wgLocal)
	}
	wgLocal.Wait()
	wg.Done()
}

// Interface that all sources should satisfy
type Source interface {
	Listen(wg sync.WaitGroup)
	init(log logger.Logger, conf config.Source, rootList *listener.Listener)
}

// Helper struct that performs common functions that most if not all Sources
// will use
type Base struct {
	config.Source

	Listener *listener.Listener
	listLock sync.Mutex

	Log logger.Logger
}

func (s *Base) init(log logger.Logger, conf config.Source, rootList *listener.Listener) {
	s.Log = log
	s.Url = conf.Url
	s.Region = conf.Region
	s.Port = conf.Port
	s.BindIp = conf.BindIp
	s.Listener = rootList
}

// Convert a raw json string into an Event struct
func (s *Base) processRaw(msg string) event.Event {

	var evt event.Event
	if err := json.Unmarshal([]byte(msg), &evt); err != nil {
		s.Log.Error("Can not turn string into event: %s", msg)
		return evt
	}
	evt.Id = generateId()
	s.Log.Event(&evt)
	s.processEvent(evt)

	return evt
}

// Process the event by running it through the tree structure of listeners
func (s *Base) processEvent(evt event.Event) {

	s.listLock.Lock()
	s.Listener.Exec(evt)
	s.listLock.Unlock()
}

// Generate an id for a newly created event
func generateId() string {

	size := 32 // change the length of the generated random string here
	rb := make([]byte, size)
	_, err := rand.Read(rb)
	if err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(rb)
}

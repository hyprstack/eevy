package main

import (
	"os"
	"sync"

	"github.com/hevnly/eevy/config"
	"github.com/hevnly/eevy/listener"
	"github.com/hevnly/eevy/source"

	"github.com/op/go-logging"
)

/**
 * Create the loggers
 */
var gLog = logging.MustGetLogger("applog")

func main() {

	backend1 := logging.NewLogBackend(os.Stderr, "", 0)
	logging.AddModuleLevel(backend1)

	c := config.Config{}
	c.LoadFromFile("./conf.yml")

	rootListener := buildListeners(&c.Listeners)

	var wg sync.WaitGroup
	wg.Add(1)
	go startSources(&c.Sources, rootListener, wg)
	wg.Wait()
}

func buildListeners(conf *config.ListenerList) *listener.EventListener {
	rootListener := listener.EventListener{}
	rootListener.Name = ""
	for evtName, listners := range *conf {
		for _, l := range listners {

			list := listener.BuildFromConf(&l)
			rootListener.Add(evtName, list)
		}
	}
	return &rootListener
}

func startSources(sourceConf *[]config.Source, rootList *listener.EventListener, wg sync.WaitGroup) {

	var wgLocal sync.WaitGroup
	var sources []source.Source
	for _, conf := range *sourceConf {
		tmp := *source.BuildFromConfig(conf, rootList)

		sources = append(sources, tmp)
		wgLocal.Add(1)
		go tmp.Listen(wgLocal)
	}
	wgLocal.Wait()
	wg.Done()
}

package main

import (
	"sync"

	"github.com/hevnly/eevy/config"
	"github.com/hevnly/eevy/handler"
	"github.com/hevnly/eevy/listener"
	"github.com/hevnly/eevy/source"
)

func main() {

	confPath := "./conf.yml"

	c := config.Config{}
	c.LoadFromFile(confPath)

	log := NewLogger(&c.Logs)
	log.Info("Reading configuration file %s", confPath)

	handlers := handler.BuildFromConf(c.Handlers, log)
	rootListener := listener.BuildListener(c.Listeners, handlers, log)

	var wg sync.WaitGroup
	wg.Add(1)
	go source.StartSources(&c.Sources, rootListener, log, wg)
	wg.Wait()
}

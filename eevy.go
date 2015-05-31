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

	log := NewLogger()

	log.Info("Reading configuration file %s", confPath)

	c := config.Config{}
	c.LoadFromFile(confPath)

	handlers := handler.BuildFromConf(c.Handlers, log)
	rootListener := listener.BuildListener(c.Listeners, handlers, log)

	var wg sync.WaitGroup
	wg.Add(1)
	go source.StartSources(&c.Sources, rootListener, log, wg)
	wg.Wait()
}

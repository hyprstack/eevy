package main

import (
	"sync"

	"github.com/hevnly/eevy/config"
	"github.com/hevnly/eevy/listener"
	"github.com/hevnly/eevy/source"
)

func main() {

	confPath := "./conf.yml"

	log := NewLogger()

	log.Info("Reading configuration file %s", confPath)

	c := config.Config{}
	c.LoadFromFile(confPath)

	rootListener := listener.BuildListeners(&c.Listeners, log)

	var wg sync.WaitGroup
	wg.Add(1)
	go source.StartSources(&c.Sources, rootListener, log, wg)
	wg.Wait()
}

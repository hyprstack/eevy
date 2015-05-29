package main

import (
	"sync"

	"github.com/hevnly/eevy/config"
	"github.com/hevnly/eevy/listener"
	"github.com/hevnly/eevy/source"
)

func main() {

	log := NewLogger()

	c := config.Config{}
	c.LoadFromFile("./conf.yml")

	rootListener := listener.BuildListeners(&c.Listeners, log)

	var wg sync.WaitGroup
	wg.Add(1)
	go source.StartSources(&c.Sources, rootListener, wg)
	wg.Wait()
}

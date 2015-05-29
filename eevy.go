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

	rootListener := listener.BuildListeners(&c.Listeners)

	var wg sync.WaitGroup
	wg.Add(1)
	go source.StartSources(&c.Sources, rootListener, wg)
	wg.Wait()
}

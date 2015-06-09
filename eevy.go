package main

import (
	"flag"
	"fmt"
	"sync"

	"github.com/hevnly/eevy/config"
	"github.com/hevnly/eevy/handler"
	"github.com/hevnly/eevy/listener"
	"github.com/hevnly/eevy/source"
)

var version string

func main() {

	var verFlag bool
	flag.BoolVar(&verFlag, "version", false, "Returns the version number")
	var confPath string
	flag.StringVar(&confPath, "conf", "conf.yml", "Location on config file (default conf.yml)")
	flag.Parse()

	// if the user wants to see the version
	if verFlag {
		fmt.Println(version)
		return
	}

	c := config.Config{}
	c.LoadFromFile(confPath)

	log := NewLogger(&c.Logs)
	log.Info("eevy starting (v%s)", version)
	log.Info("Reading configuration file %s", confPath)

	handlers := handler.BuildFromConf(c.Handlers, log)
	rootListener := listener.BuildListener(c.Listeners, handlers, log)

	var wg sync.WaitGroup
	wg.Add(1)
	go source.StartSources(&c.Sources, rootListener, log, wg)
	wg.Wait()
}

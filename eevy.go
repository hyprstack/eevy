package main

import (
	"flag"
	"fmt"
	"os"
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
	flag.StringVar(&confPath, "conf", "conf.yml", "Location on config file")
	flag.Parse()

	// if the user wants to see the version
	if verFlag {
		fmt.Printf("%s", version)
		os.Exit(0)
	}

	c := config.Config{}
	err := c.LoadFromPath(confPath)
	if err != nil {
		fmt.Printf("Error loading config file: %s\n\n", err.Error())
		os.Exit(2)
	}

	log := NewLogger(&c.Logs)
	log.Info("eevy starting (v%s)", version)
	log.Info("Reading configuration file %s", confPath)

	log.Info("Building handlers")
	handlers := handler.BuildFromConf(c.Handlers, log)
	log.Info("Building listeners")
	rootListener := listener.BuildListener(c.Listeners, handlers, log)

	log.Info("Starting sources")
	var wg sync.WaitGroup
	wg.Add(1)
	go source.StartSources(&c.Sources, rootListener, log, wg)
	wg.Wait()
}

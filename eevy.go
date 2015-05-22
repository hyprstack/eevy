package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"hevnly/eevy/config"
	"hevnly/eevy/listener"
	"hevnly/eevy/source"

	"github.com/op/go-logging"
	"gopkg.in/yaml.v2"
)

/**
 * Create the loggers
 */
var gLog = logging.MustGetLogger("applog")

func init() {
	backend1 := logging.NewLogBackend(os.Stderr, "", 0)
	logging.AddModuleLevel(backend1)
}

/**
 * Read the configuration file
 */
var gConfig = config.Config{}
var gSrc = source.Sqs{}

func init() {

	gLog.Info("Reading config file './conf.yml'")

	filename, _ := filepath.Abs("./conf.yml")
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(yamlFile, &gConfig)
	if err != nil {
		panic(err)
	}
}

func main() {

	rootListener := buildListeners(&gConfig.Listeners)

	var wg sync.WaitGroup
	wg.Add(1)
	go startSources(&gConfig.Sources, rootListener, wg)
	wg.Wait()
}

func buildListeners(conf *map[string][]config.Listener) *listener.EventListener {
	rootListener := listener.EventListener{}
	rootListener.Name = ""
	for evtName, listners := range gConfig.Listeners {
		for _, listConf := range listners {

			list := listener.BuildFromConf(listConf)
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

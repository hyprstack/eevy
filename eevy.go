package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"github.com/hevnly/eevy/config"
	"github.com/hevnly/eevy/listener"
	"github.com/hevnly/eevy/source"

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

// Various stucts that are used to configure sources and handlers.
package config

import (
	"fmt"
	"os"

	"encoding/json"
	"io/ioutil"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// Represents the structure that should be used in configuration files
type Config struct {
	Logs      Logger
	Sources   []Source
	Handlers  HandlerList
	Listeners ListenerList
}

type ListenerList map[string][]string

type HandlerList map[string]Handler

type Handler map[string]interface{}

func (this *Handler) get(s string) string {

	if _, ok := (*this)[s]; !ok {
		return ""
	}
	if t, ok := (*this)[s].(string); ok {
		return t
	}
	return ""
}

// Return the type of handler eg "sqs", "lambda" etc
func (this *Handler) GetType() string {
	return this.get("type")
}

// Return the message that should be sent by the handler, default is "${message}"
func (this *Handler) GetMessage() string {

	r := this.get("message")
	if r != "" {
		return r
	}
	return "${message}"
}

//  Convert this handler into a json string
func (this *Handler) String() string {
	b, _ := json.Marshal(this)
	return string(b)
}

func (this *Handler) Init(s string) {
	return
}

func (this *Config) LoadFromPath(s string) error {

	fi, err := os.Stat(s)
	if err != nil {
		return err
	}
	fmt.Printf("Reading configuration from %s\n", s)
	if fi.IsDir() {
		this.loadFromFolder(s)
	} else {
		this.loadFromFile(s)
	}
	return nil
}

func (this *Config) loadFromFolder(s string) {

	fileList := []string{}
	filepath.Walk(s, func(path string, f os.FileInfo, err error) error {

		if !f.IsDir() {
			fileList = append(fileList, path)
		}
		return nil
	})

	for _, file := range fileList {

		tmp := Config{}
		tmp.loadFromFile(file)
		this.merge(&tmp)
	}
}

// Build this configuration struct from a configuration file
func (this *Config) loadFromFile(s string) {

	filename, _ := filepath.Abs(s)
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(yamlFile, &this)
	if err != nil {
		panic(err)
	}
}

func (this *Config) merge(src *Config) {

	this.mergeHandlers(&src.Handlers)
	this.mergeSources(&src.Sources)
	this.mergeListeners(&src.Listeners)
	this.mergeLogs(&src.Logs)
}

// type HandlerList map[string]Handler
// type Handler map[string]interface{}

func (this *Config) mergeHandlers(src *HandlerList) {

	if this.Handlers == nil {
		this.Handlers = make(HandlerList)
	}
	for k, v := range *src {

		this.Handlers[k] = v
	}
}

func (this *Config) mergeSources(src *[]Source) {

	this.Sources = append(this.Sources, *src...)
}

func (this *Config) mergeListeners(src *ListenerList) {

	if this.Listeners == nil {
		this.Listeners = make(ListenerList)
	}
	for k, v := range *src {

		if val, ok := this.Listeners[k]; ok {
			this.Listeners[k] = append(this.Listeners[k], val...)
		} else {
			this.Listeners[k] = v
		}
	}
}

func (this *Config) mergeLogs(src *Logger) {

	if src.Event != "" {
		this.Logs.Event = src.Event
	}
	if src.App != "" {
		this.Logs.App = src.App
	}
	if src.Handler != "" {
		this.Logs.Handler = src.Handler
	}
	if src.Level != "" {
		this.Logs.Level = src.Level
	}
}

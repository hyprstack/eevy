// handler describes all of the systems that events can be sent to by eevy
package handler

import (
	"regexp"
	"strings"

	"encoding/json"

	appConfig "github.com/hevnly/eevy/config"
	appHandConfig "github.com/hevnly/eevy/config/handler"
	"github.com/hevnly/eevy/event"
	"github.com/hevnly/eevy/handler/config"
	"github.com/hevnly/eevy/logger"
)

type HandlerList struct {
	List map[string]Handler
}

type Handler interface {
	Exec(evt event.Event)
	GetType() string
	GetName() string
	SetName(s string)
	GetConfig() config.Handler
}

type HandlerBase struct {
	Log  logger.Logger
	Name string
}

func (this *HandlerBase) GetName() string {
	return this.Name
}

func (this *HandlerBase) SetName(s string) {
	this.Name = s
}

func BuildFromConf(conf appConfig.HandlerList, log logger.Logger) *HandlerList {
	hl := &HandlerList{
		List: make(map[string]Handler),
	}
	for name, c := range conf {
		h := *BuildHandlerFromConf(&c, log)
		h.SetName(name)

		log.Debug("Created handler, %s, of type %s", h.GetName(), h.GetType())
		hl.List[name] = h
	}
	return hl
}

func (this *HandlerList) Get(s string) Handler {
	if h, ok := this.List[s]; ok {
		return h
	}
	return nil
}

// Receives a configuration struct and creates the relevant Handler
func BuildHandlerFromConf(conf config.Handler, log logger.Logger) *Handler {

	var l Handler
	switch conf.GetType() {
	case "sqs":
		tl := &Sqs{Config: &appHandConfig.Sqs{}}
		tl.Log = log
		tl.Config.Init(conf.String())
		l = tl

	case "lambda":
		tl := &Lambda{Config: &appHandConfig.Lambda{}}
		tl.Log = log
		tl.Config.Init(conf.String())
		l = tl

	case "oauth2":
		tl := &OAuth2{Config: &appHandConfig.OAuth2{}}
		tl.Log = log
		tl.Config.Init(conf.String())
		l = tl

	case "cli":
		tl := &Cli{Config: &appHandConfig.Cli{}}
		tl.Log = log
		tl.Config.Init(conf.String())
		l = tl

	default:
		log.Warning("Could not create handler for type '%s'", conf.GetType())
		return nil

	}
	return &l
}

// Replaces variables ("${}") in the string to their actual value
func magicString(s string, evt event.Event) string {

	rep := regexp.MustCompile("(\\${|})")
	rst := findMagicStrings(s)
	for _, v := range rst {
		variable := rep.ReplaceAllString(v, "")
		opt := strings.Split(variable, ".")
		switch opt[0] {
		case "message":
			str := ""
			if len(opt) <= 1 {
				b, _ := json.Marshal(evt.Message)
				str = string(b)
			} else {
				str = evt.GetString(strings.Join(opt[1:], "."))
			}
			s = strings.Replace(s, v, str, -1)
		}
	}
	return s
}

// Finds all of the ${} in the given string
func findMagicStrings(s string) []string {
	re := regexp.MustCompile("\\${(.*?)}")
	return re.FindAllString(s, -1)
}

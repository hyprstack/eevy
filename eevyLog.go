package main

import (
	"os"

	"github.com/op/go-logging"

	"github.com/hevnly/eevy/logger"
)

type EevyLog struct {
	AppLog      logging.Logger
	EventLog    logging.Logger
	ListenerLog logging.Logger
}

func NewLogger() *EevyLog {
	log := EevyLog{
		AppLog:      logging.Logger{Module: "appLog"},
		EventLog:    logging.Logger{Module: "eventLog"},
		ListenerLog: logging.Logger{Module: "listnerLog"},
	}
	backend1 := logging.NewLogBackend(os.Stderr, "", 0)
	backend1Leveled := logging.AddModuleLevel(backend1)
	backend1Leveled.SetLevel(logging.ERROR, "")

	log.AppLog.SetBackend(backend1Leveled)

	return &log
}

func (this *EevyLog) Event(evt logger.Event) {
	this.Critical("%s %s", evt.GetName(), evt.GetId())
}

func (this *EevyLog) Listener(l logger.Listener, e logger.Event) {
	this.Critical("%s %s %s", l.GetType(), e.GetName(), e.GetId())
}

func (this *EevyLog) ListenerError(l logger.Listener, e logger.Event) {

}

func (this *EevyLog) Critical(format string, args ...interface{}) {
	this.AppLog.Critical(format, args...)
}

func (this *EevyLog) Error(format string, args ...interface{}) {
	this.AppLog.Critical(format, args...)
}

func (this *EevyLog) Warning(format string, args ...interface{}) {
	this.AppLog.Critical(format, args...)
}

func (this *EevyLog) Notice(format string, args ...interface{}) {
	this.AppLog.Critical(format, args...)
}

func (this *EevyLog) Info(format string, args ...interface{}) {
	this.AppLog.Critical(format, args...)
}

func (this *EevyLog) Debug(format string, args ...interface{}) {
	this.AppLog.Critical(format, args...)
}

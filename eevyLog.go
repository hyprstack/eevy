package main

import (
	"os"

	"github.com/op/go-logging"

	"github.com/hevnly/eevy/logger"
)

type EevyLog struct {
	AppLog     logging.Logger
	EventLog   logging.Logger
	HandlerLog logging.Logger
}

func NewLogger() *EevyLog {

	log := EevyLog{
		AppLog:     logging.Logger{Module: "appLog"},
		EventLog:   logging.Logger{Module: "eventLog"},
		HandlerLog: logging.Logger{Module: "handlerLog"},
	}
	log.buildAppLog()
	log.buildEventLog()
	log.buildHandlerLog()

	return &log
}

func (this *EevyLog) buildAppLog() {

	fo, err := os.Create("/var/log/eevy/app.log")
	if err != nil {
		return
	}
	appBe := logging.NewLogBackend(fo, "", 0)
	var appFormat = logging.MustStringFormatter(
		"%{time} %{level} %{message}",
	)
	appBeFormatter := logging.NewBackendFormatter(appBe, appFormat)
	appLeveled := logging.AddModuleLevel(appBeFormatter)
	appLeveled.SetLevel(logging.DEBUG, "")
	this.AppLog.SetBackend(appLeveled)
}

func (this *EevyLog) buildEventLog() {

	fo, err := os.Create("/var/log/eevy/event.log")
	if err != nil {
		return
	}
	evtBe := logging.NewLogBackend(fo, "", 0)
	var evtFormat = logging.MustStringFormatter(
		"%{time} %{message}",
	)
	evtBeFormatter := logging.NewBackendFormatter(evtBe, evtFormat)
	evtLeveled := logging.AddModuleLevel(evtBeFormatter)
	evtLeveled.SetLevel(logging.DEBUG, "")
	this.EventLog.SetBackend(evtLeveled)
}

func (this *EevyLog) buildHandlerLog() {

	fo, err := os.Create("/var/log/eevy/handler.log")
	if err != nil {
		return
	}
	handBe := logging.NewLogBackend(fo, "", 0)
	var handFormat = logging.MustStringFormatter(
		"%{time} %{message}",
	)
	handBeFormatter := logging.NewBackendFormatter(handBe, handFormat)
	handLeveled := logging.AddModuleLevel(handBeFormatter)
	handLeveled.SetLevel(logging.DEBUG, "")
	this.HandlerLog.SetBackend(handLeveled)
}

func (this *EevyLog) Event(evt logger.Event) {
	this.EventLog.Critical("%s %s", evt.GetName(), evt.GetId())
}

func (this *EevyLog) Handler(l logger.Handler, e logger.Event) {
	this.HandlerLog.Info("EXEC %s %s %s", l.GetName(), e.GetName(), e.GetId())
}

func (this *EevyLog) HandlerError(l logger.Handler, e logger.Event) {
	this.HandlerLog.Error("ERROR %s %s %s", l.GetName(), e.GetName(), e.GetId())
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

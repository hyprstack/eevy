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
	log.buildAppLog()
	log.buildEventLog()
	log.buildListenerLog()

	return &log
}

func (this *EevyLog) buildAppLog() {

	appBe := logging.NewLogBackend(os.Stdout, "", 0)
	var appFormat = logging.MustStringFormatter(
		"%{time} %{level} %{message}",
	)
	appBeFormatter := logging.NewBackendFormatter(appBe, appFormat)
	appLeveled := logging.AddModuleLevel(appBeFormatter)
	appLeveled.SetLevel(logging.DEBUG, "")
	this.AppLog.SetBackend(appLeveled)
}

func (this *EevyLog) buildEventLog() {

	evtBe := logging.NewLogBackend(os.Stdout, "", 0)
	var evtFormat = logging.MustStringFormatter(
		"%{time} %{message}",
	)
	evtBeFormatter := logging.NewBackendFormatter(evtBe, evtFormat)
	evtLeveled := logging.AddModuleLevel(evtBeFormatter)
	evtLeveled.SetLevel(logging.DEBUG, "")
	this.EventLog.SetBackend(evtLeveled)
}

func (this *EevyLog) buildListenerLog() {

	listBe := logging.NewLogBackend(os.Stdout, "", 0)
	var listFormat = logging.MustStringFormatter(
		"%{time} %{message}",
	)
	listBeFormatter := logging.NewBackendFormatter(listBe, listFormat)
	listLeveled := logging.AddModuleLevel(listBeFormatter)
	listLeveled.SetLevel(logging.DEBUG, "")
	this.ListenerLog.SetBackend(listLeveled)
}

func (this *EevyLog) Event(evt logger.Event) {
	this.EventLog.Critical("%s %s", evt.GetName(), evt.GetId())
}

func (this *EevyLog) Listener(l logger.Listener, e logger.Event) {
	this.ListenerLog.Info("EXEC %s %s %s", l.GetType(), e.GetName(), e.GetId())
}

func (this *EevyLog) ListenerError(l logger.Listener, e logger.Event) {
	this.ListenerLog.Error("ERROR %s %s %s", l.GetType(), e.GetName(), e.GetId())
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

package main

import (
	"fmt"
	"io"
	"os"

	"github.com/op/go-logging"

	"github.com/hevnly/eevy/logger"
)

type EevyLog struct {
	AppLog     logging.Logger
	EventLog   logging.Logger
	HandlerLog logging.Logger
	Config     EevyLogConfig
}

type EevyLogConfig interface {
	GetEventPath() string
	GetHandlerPath() string
	GetAppPath() string
	GetSeverityLevel() string
}

func NewLogger(config EevyLogConfig) *EevyLog {

	log := EevyLog{
		AppLog:     logging.Logger{Module: "appLog"},
		EventLog:   logging.Logger{Module: "eventLog"},
		HandlerLog: logging.Logger{Module: "handlerLog"},
		Config:     config,
	}
	log.buildAppLog()
	log.buildEventLog()
	log.buildHandlerLog()

	return &log
}

func getWriter(s string) io.Writer {

	var fo io.Writer
	var err error
	switch s {
	case "stdout":
		fo = os.Stdout

	case "stderr":
		fo = os.Stderr

	default:
		fo, err = os.OpenFile(s, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0660)
	}
	if err != nil {
		fmt.Printf("Error creating log file: %s\n", err.Error())
		return nil
	}
	return fo
}

func (this *EevyLog) buildAppLog() {

	fo := getWriter(this.Config.GetAppPath())
	if fo == nil {
		return
	}
	appBe := logging.NewLogBackend(fo, "", 0)
	var appFormat = logging.MustStringFormatter(
		"%{time} APP %{level} %{message}",
	)
	appBeFormatter := logging.NewBackendFormatter(appBe, appFormat)
	appLeveled := logging.AddModuleLevel(appBeFormatter)
	level, _ := logging.LogLevel(this.Config.GetSeverityLevel())

	appLeveled.SetLevel(level, "")
	this.AppLog.SetBackend(appLeveled)

	this.Debug("App log created: %s", this.Config.GetAppPath())
}

func (this *EevyLog) buildEventLog() {

	fo := getWriter(this.Config.GetEventPath())
	if fo == nil {
		return
	}
	evtBe := logging.NewLogBackend(fo, "", 0)
	var evtFormat = logging.MustStringFormatter(
		"%{time} EVT %{message}",
	)
	evtBeFormatter := logging.NewBackendFormatter(evtBe, evtFormat)
	evtLeveled := logging.AddModuleLevel(evtBeFormatter)
	evtLeveled.SetLevel(logging.DEBUG, "")
	this.EventLog.SetBackend(evtLeveled)

	this.Debug("Event log created: %s", this.Config.GetEventPath())
}

func (this *EevyLog) buildHandlerLog() {

	fo := getWriter(this.Config.GetHandlerPath())
	if fo == nil {
		return
	}
	handBe := logging.NewLogBackend(fo, "", 0)
	var handFormat = logging.MustStringFormatter(
		"%{time} LIST %{message}",
	)
	handBeFormatter := logging.NewBackendFormatter(handBe, handFormat)
	handLeveled := logging.AddModuleLevel(handBeFormatter)
	handLeveled.SetLevel(logging.DEBUG, "")
	this.HandlerLog.SetBackend(handLeveled)

	this.Debug("Handler log created: %s", this.Config.GetEventPath())
}

func (this *EevyLog) Event(evt logger.Event) {
	this.EventLog.Critical("%s %s", evt.GetName(), evt.GetId())
}

func (this *EevyLog) Handler(l logger.Handler, e logger.Event) {
	this.HandlerLog.Info("EXEC %s %s %s", l.GetName(), e.GetName(), e.GetId())
}

func (this *EevyLog) HandlerError(l logger.Handler, msg string, e logger.Event) {
	this.HandlerLog.Error("ERROR %s %s %s %s", l.GetName(), e.GetName(), e.GetId(), msg)
}

func (this *EevyLog) Critical(format string, args ...interface{}) {
	this.AppLog.Critical(format, args...)
}

func (this *EevyLog) Error(format string, args ...interface{}) {
	this.AppLog.Error(format, args...)
}

func (this *EevyLog) Warning(format string, args ...interface{}) {
	this.AppLog.Warning(format, args...)
}

func (this *EevyLog) Notice(format string, args ...interface{}) {
	this.AppLog.Notice(format, args...)
}

func (this *EevyLog) Info(format string, args ...interface{}) {
	this.AppLog.Info(format, args...)
}

func (this *EevyLog) Debug(format string, args ...interface{}) {
	this.AppLog.Debug(format, args...)
}

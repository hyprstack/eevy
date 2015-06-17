package mocks

import "github.com/hevnly/eevy/logger"
import "github.com/stretchr/testify/mock"

type Logger struct {
	mock.Mock
}

func (m *Logger) Critical(format string, args ...interface{}) {
	m.Called(format, args)
}
func (m *Logger) Error(format string, args ...interface{}) {
	m.Called(format, args)
}
func (m *Logger) Warning(format string, args ...interface{}) {
	m.Called(format, args)
}
func (m *Logger) Notice(format string, args ...interface{}) {
	m.Called(format, args)
}
func (m *Logger) Info(format string, args ...interface{}) {
	m.Called(format, args)
}
func (m *Logger) Debug(format string, args ...interface{}) {
	m.Called(format, args)
}
func (m *Logger) Event(evt logger.Event) {
	m.Called(evt)
}
func (m *Logger) Handler(l logger.Handler, e logger.Event) {
	m.Called(l, e)
}
func (m *Logger) HandlerError(l logger.Handler, msg string, e logger.Event) {
	m.Called(l, msg, e)
}

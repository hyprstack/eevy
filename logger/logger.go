package logger

type Event interface {
	GetName() string
	GetId() string
}

type Listener interface {
	GetType() string
}

type Logger interface {

	// Critical logs a message using CRITICAL as log level.
	Critical(format string, args ...interface{})

	// Error logs a message using ERROR as log level.
	Error(format string, args ...interface{})

	// Warning logs a message using WARNING as log level.
	Warning(format string, args ...interface{})

	// Notice logs a message using NOTICE as log level.
	Notice(format string, args ...interface{})

	// Info logs a message using INFO as log level.
	Info(format string, args ...interface{})

	// Debug logs a message using DEBUG as log level.
	Debug(format string, args ...interface{})

	// Logs an event being created
	Event(evt Event)

	// logs a listener being triggered
	Listener(l Listener, e Event)

	// logs a listener error
	ListenerError(l Listener, e Event)
}

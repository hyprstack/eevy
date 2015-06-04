package config

type Logger struct {
	Event   string
	Handler string
	App     string
	Level   string
}

func (this *Logger) GetEventPath() string {
	return this.Event
}

func (this *Logger) GetHandlerPath() string {
	return this.Handler
}

func (this *Logger) GetAppPath() string {
	return this.App
}

func (this *Logger) GetSeverityLevel() string {
	return this.Level
}

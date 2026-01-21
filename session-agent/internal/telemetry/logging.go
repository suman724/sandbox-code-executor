package telemetry

import "log"

type Logger struct {
	component string
}

func NewLogger(component string) Logger {
	return Logger{component: component}
}

func (l Logger) Info(msg string) {
	log.Printf("component=%s level=info msg=%q", l.component, msg)
}

func (l Logger) Error(msg string) {
	log.Printf("component=%s level=error msg=%q", l.component, msg)
}

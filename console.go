package godom

import (
	"fmt"
	"syscall/js"
)

type ConsoleLogger struct {
	console js.Value
}

var Console = &ConsoleLogger{
	console: global.Get("console"),
}

func (c *ConsoleLogger) write(level string, msg string, args ...any) {
	c.console.Call(level, fmt.Sprintf(msg, args...))
}

func (c *ConsoleLogger) Log(msg string, args ...any) {
	c.write("log", fmt.Sprintf(msg, args...))
}

func (c *ConsoleLogger) Debug(msg string, args ...any) {
	c.write("debug", fmt.Sprintf(msg, args...))
}

func (c *ConsoleLogger) Info(msg string, args ...any) {
	c.write("info", fmt.Sprintf(msg, args...))
}

func (c *ConsoleLogger) Warn(msg string, args ...any) {
	c.write("warn", fmt.Sprintf(msg, args...))
}

func (c *ConsoleLogger) Error(msg string, args ...any) {
	c.write("error", fmt.Sprintf(msg, args...))
}

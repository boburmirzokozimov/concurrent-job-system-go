package logger

import (
	"fmt"
	"os"
	"sync"
	"time"
)

type ConsoleLogger struct {
	mu sync.Mutex
}

func NewConsoleLogger() *ConsoleLogger {
	return &ConsoleLogger{}
}

func (c *ConsoleLogger) log(level string, color string, msg string, args ...any) {
	c.mu.Lock()
	defer c.mu.Unlock()

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	formatted := fmt.Sprintf(msg, args...)
	_, err := fmt.Fprintf(os.Stdout, "%s[%s] [%s] %s\033[0m\n", color, timestamp, level, formatted)
	if err != nil {
		return
	}
}

func (c *ConsoleLogger) Info(msg string, args ...any) {
	c.log("INFO", "\033[34m", msg, args...)
}

func (c *ConsoleLogger) Error(msg string, args ...any) {
	c.log("ERROR", "\033[31m", msg, args...)
}

func (c *ConsoleLogger) Warn(msg string, args ...any) {
	c.log("WARN", "\033[33m", msg, args...)
}

func (c *ConsoleLogger) Debug(msg string, args ...any) {
	c.log("DEBUG", "\033[36m", msg, args...)
}

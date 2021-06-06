package strategy

import (
	"fmt"
	"strings"
)

// logger is a simple in-memory logger.
type logger struct {
	prefix   string
	messages []string
}

// Log formats and logs a message.
func (l *logger) Log(format string, a ...interface{}) {
	l.messages = append(l.messages, fmt.Sprintf(l.prefix+format, a...))
}

// String implements fmt.Stringer.
func (l *logger) String() string {
	return strings.Join(l.messages, "\n")
}

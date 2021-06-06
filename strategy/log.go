package strategy

import (
	"fmt"
	"strings"
)

// logger is a simple in-memory logger.
type logger struct {
	prefix string
	sb     strings.Builder
}

// Log formats and logs a message.
func (l *logger) Log(format string, a ...interface{}) {
	l.sb.WriteString(fmt.Sprintf(l.prefix+format, a...))
}

// String implements fmt.Stringer.
func (l *logger) String() string {
	return l.sb.String()
}

package logger

import (
	"fmt"
	"os"
)

const (
	QUIET LogLevel = iota + 1
	ERROR
	WARN
	INFO
	DEBUG
)

type LogLevel int

var Level = WARN

func Debugf(format string, a ...any) {
	if Level >= DEBUG {
		fmt.Fprintf(os.Stderr, "DEBUG: "+format+"\n", a...)
	}
}

func Infof(format string, a ...any) {
	if Level >= INFO {
		fmt.Fprintf(os.Stderr, "INFO: "+format+"\n", a...)
	}
}

func Warnf(format string, a ...any) {
	if Level >= WARN {
		fmt.Fprintf(os.Stderr, "WARN: "+format+"\n", a...)
	}
}

func Errorf(format string, a ...any) {
	if Level >= ERROR {
		fmt.Fprintf(os.Stderr, "ERROR: "+format+"\n", a...)
	}
}

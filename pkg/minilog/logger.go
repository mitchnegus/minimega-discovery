// Copyright 2018 National Technology & Engineering Solutions of Sandia, LLC
// (NTESS). Under the terms of Contract DE-NA0003525 with NTESS, the U.S.
// Government retains certain rights in this software.

package minilog

import (
	"fmt"
	golog "log"
	"runtime"
	"strconv"
	"strings"
)

type minilogger struct {
	*golog.Logger
	Level   Level
	Color   bool // print in color
	filters []string
}

func (l *minilogger) prologue(level Level, name string) (msg string) {
	switch level {
	case DEBUG:
		msg += "DEBUG "
	case INFO:
		msg += "INFO "
	case WARN:
		msg += "WARN "
	case ERROR:
		msg += "ERROR "
	default:
		msg += "FATAL "
	}

	if name == "" {
		_, file, line, _ := runtime.Caller(4)
		short := file
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				short = file[i+1:]
				break
			}
		}
		msg += short + ":" + strconv.Itoa(line) + ": "
	} else {
		msg += name + ": "
	}

	if l.Color {
		msg = colorLine + msg
		switch level {
		case DEBUG:
			msg += colorDebug
		case INFO:
			msg += colorInfo
		case WARN:
			msg += colorWarn
		case ERROR:
			msg += colorError
		default:
			msg += colorFatal
		}
	}
	return
}

func (l *minilogger) epilogue() string {
	if l.Color {
		return Reset
	}
	return ""
}

func (l *minilogger) log(level Level, name, format string, arg ...interface{}) {
	msg := l.prologue(level, name) + fmt.Sprintf(format, arg...) + l.epilogue()
	for _, f := range l.filters {
		if strings.Contains(msg, f) {
			return
		}
	}
	l.Print(msg)
}

func (l *minilogger) logln(level Level, name string, arg ...interface{}) {
	msg := l.prologue(level, name) + fmt.Sprint(arg...) + l.epilogue()
	for _, f := range l.filters {
		if strings.Contains(msg, f) {
			return
		}
	}
	l.Println(msg)
}

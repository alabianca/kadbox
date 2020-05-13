package log

import (
	"log"
	"os"
)

var debug *log.Logger
var error *log.Logger
var info *log.Logger

func init() {
	debug = log.New(os.Stdout, "[ debug ] ", log.Ltime)
	error = log.New(os.Stderr, "[ error ] ", log.Ltime)
	info = log.New(os.Stdout, "[ info ] ", log.Ltime)
}

func Info(v ...interface{}) {
	info.Println(v...)
}

func Infof(format string, v ...interface{}) {
	info.Printf(format, v...)
}

func Error(v ...interface{}) {
	error.Println(v...)
}

func Errorf(format string, v ...interface{}) {
	error.Printf(format, v...)
}

func Debug(v ...interface{}) {
	debug.Println(v...)
}

func Debugf(format string, v ...interface{}) {
	debug.Printf(format, v...)
}

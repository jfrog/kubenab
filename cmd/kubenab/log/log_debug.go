// +build !strip_debug

package log

import "log"

func Debug(v ...interface{}) {
	log.Print(v...)
}

func Debugln(v ...interface{}) {
	log.Println(v...)
}

func Debugf(fmt string, args ...interface{}) {
	log.Printf(fmt, args...)
}

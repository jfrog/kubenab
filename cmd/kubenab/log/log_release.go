// +build strip_debug,!debug

package log

func Debug(v ...interface{}) {}

func Debugln(v ...interface{}) {}

func Debugf(fmt string, args ...interface{}) {}

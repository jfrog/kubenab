package log

import _log "github.com/sirupsen/logrus"

func Panicf(fmt string, v ...interface{}) {
	_log.Panicf(fmt, v...)
}

func Panicln(v ...interface{}) {
	_log.Panicln(v...)
}

func Fatalf(fmt string, v ...interface{}) {
	_log.Fatalf(fmt, v...)
}

func Fatalln(v ...interface{}) {
	_log.Fatalln(v...)
}

func Errorf(fmt string, v ...interface{}) {
	_log.Errorf(fmt, v...)
}

func Errorln(v ...interface{}) {
	_log.Fatalln(v...)
}

func Warnln(v ...interface{}) {
	_log.Println(v...)
}

func Warnf(fmt string, v ...interface{}) {
	_log.Warnf(fmt, v...)
}

func Println(v ...interface{}) {
	_log.Println(v...)
}

func Printf(fmt string, v ...interface{}) {
	_log.Printf(fmt, v...)
}

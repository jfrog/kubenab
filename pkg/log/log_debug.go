// +build debug

package log

import _log "github.com/sirupsen/logrus"

func Debugf(fmt string, v ...interface{}) {
	_log.Debugf(fmt, v...)
}

func Debugln(v ...interface{}) {
	_log.Debugln(v...)
}

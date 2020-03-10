// log package is a wrapper around the current log package to provide the
// possibility to "disable" Debug log output => GoLang compiler will
// remove it from the binary.
package log

import "log"

func Fatal(v ...interface{}) {
	log.Fatal(v...)
}

func Fatalln(v ...interface{}) {
	log.Fatalln(v...)
}

func Println(v ...interface{}) {
	log.Println(v...)
}

func Printf(fmt string, args ...interface{}) {
	log.Printf(fmt, args...)
}

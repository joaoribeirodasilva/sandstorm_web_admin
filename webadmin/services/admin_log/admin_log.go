package admin_log

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

type Severity uint8

const (
	LOG_CRITICAL Severity = iota
	LOG_ERROR
	LOG_WARNING
	LOG_INFO
	LOG_DEBUG
)

type Log struct {
	path   string
	file   *os.File
	out    *os.File
	level  Severity
	multi  io.Writer
	isOpen bool
}

var severities = map[Severity]string{
	LOG_CRITICAL: "CRITICAL",
	LOG_ERROR:    "ERROR",
	LOG_WARNING:  "WARNING",
	LOG_INFO:     "INFO",
	LOG_DEBUG:    "DEBUG",
}

func New() *Log {

	var err error
	l := new(Log)
	now := time.Now()
	fName := fmt.Sprintf("./logs/%04d_%02d_%02dT%02d_%02d_%02d.log", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())
	l.path, err = filepath.Abs(fName)
	if err != nil {
		log.Fatalf("failed to calculate absolute file for path '%s'. ERR:", fName)
	}
	l.out = os.Stdout
	l.level = LOG_INFO

	return l
}

func (l *Log) Open() {
	l.isOpen = true
}

func (l *Log) SetLogLevel(severity Severity) {

	l.level = severity
}

func (l *Log) Write(msg string, module string, severity Severity) error {

	if l.level <= severity {

		var err error

		if l.isOpen {
			l.file, err = os.OpenFile(l.path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
			if err != nil {
				log.Fatalf("error opening file '%s': %v", l.path, err)
			}
			defer l.file.Close()
		}

		l.multi = io.MultiWriter(l.out, l.file)

		logLine := []byte(fmt.Sprintf("[%-8s] [%-10s]: %s\n", severities[severity], module, msg))
		if l.multi != nil {
			l.multi.Write(logLine)
			if severity == LOG_CRITICAL {
				os.Exit(1)
			}
		} else {
			log.Printf("%s", logLine)
		}

	}
	return fmt.Errorf(msg)

}

func (l *Log) Close() {

	if l.file != nil {
		l.file.Close()
		l.file = nil
		l.isOpen = false
	}

	if l.out != nil {
		l.out = nil
	}
}

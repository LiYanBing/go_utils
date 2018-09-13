package glog

import (
	"fmt"

	"github.com/golang/glog"
)

type Logger struct {
	reqId  string
	format string
}

func New(reqId string) *Logger {
	format := fmt.Sprintf("[%v]", reqId)
	return &Logger{
		reqId:  reqId,
		format: format,
	}
}

func (l *Logger) Info(args ...interface{}) {
	glog.InfoDepth(1, append([]interface{}{l.format}, args...)...)
}

func (l *Logger) Infof(format string, args ...interface{}) {
	glog.InfoDepth(1, append([]interface{}{l.format}, fmt.Sprintf(format, args...))...)
}

func (l *Logger) Warning(args ...interface{}) {
	glog.WarningDepth(1, append([]interface{}{l.format}, args...)...)
}

func (l *Logger) Warningf(format string, args ...interface{}) {
	glog.WarningDepth(1, append([]interface{}{l.format}, fmt.Sprintf(format, args...))...)
}

func (l *Logger) Error(args ...interface{}) {
	//TODO if argc contains [Alert] send mail
	glog.ErrorDepth(1, append([]interface{}{l.format}, args...)...)
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	//TODO if argc contains [Alert] send mail
	glog.ErrorDepth(1, append([]interface{}{l.format}, fmt.Sprintf(format, args...))...)
}

func (l *Logger) Alert(args ...interface{}) {
	l.Error(append([]interface{}{"[Alert]"}, args...)...)
}

func (l *Logger) Alertf(format string, args ...interface{}) {
	l.Errorf("[Alert] "+format, args...)
}

func (l *Logger) Fatal(args ...interface{}) {
	glog.FatalDepth(1, append([]interface{}{l.format}, args...)...)
}

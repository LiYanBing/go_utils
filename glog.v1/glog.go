package glog

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"study/stu/go_utils/mail"

	"github.com/golang/glog"
	"study/stu/go_utils/reqid.v1"
)

const (
	mailFlag = "[Alert]"
	key      = 1
)

//如果需要发送邮件需要设置下面的值
var (
	Subject string
	Mail    *mail.Mail
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

func newWithContext(ctx context.Context) *Logger {
	if reqId, ok := reqid.FromContext(ctx); ok {
		return New(reqId)
	}
	glog.V(1).Info("cannot get reqid")
	return New("")
}

func FromContext(ctx context.Context) *Logger {
	if dl, ok := ctx.Value(key).(*Logger); ok {
		return dl
	}
	return newWithContext(ctx)
}

func NewContext(ctx context.Context) context.Context {
	dl := FromContext(ctx)
	return context.WithValue(ctx, key, dl)
}

func GlogHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := NewContext(req.Context())
		next.ServeHTTP(w, req.WithContext(ctx))
	})
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
	glog.ErrorDepth(1, append([]interface{}{l.format}, args...)...)
	l.sendMail(append([]interface{}{l.format}, args...)...)
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	glog.ErrorDepth(1, append([]interface{}{l.format}, fmt.Sprintf(format, args...))...)
}

func (l *Logger) Alert(args ...interface{}) {
	l.Error(append([]interface{}{mailFlag}, args...)...)
}

func (l *Logger) Alertf(format string, args ...interface{}) {
	l.Error(append([]interface{}{mailFlag}, fmt.Sprintf(format, args...))...)
}

func (l *Logger) Fatal(args ...interface{}) {
	glog.FatalDepth(1, append([]interface{}{l.format}, args...)...)
}

var verbosePool = sync.Pool{
	New: func() interface{} {
		return new(Verbose)
	},
}

type Verbose struct {
	format string
	v      glog.Verbose
}

func (l *Logger) V(level glog.Level) *Verbose {
	verbose := verbosePool.Get().(*Verbose)
	verbose.format = l.format
	verbose.v = glog.V(level)
	return verbose
}

func (v *Verbose) Info(args ...interface{}) {
	defer verbosePool.Put(v)
	if v.v {
		glog.InfoDepth(1, v.format, args)
	}
}

func (v *Verbose) Infof(format string, args ...interface{}) {
	defer verbosePool.Put(v)
	if v.v {
		glog.InfoDepth(1, v.format, fmt.Sprintf(format, args...))
	}
}

func (l *Logger) GetReqId() string {
	return l.reqId
}

func (l *Logger) sendMail(args ...interface{}) {
	go func() {
		if Mail != nil {
			for _, v := range args {
				if str, ok := v.(string); ok && strings.Contains(str, mailFlag) {
					entity := mail.TextEntity(fmt.Sprint(args...))
					err := Mail.Send(Subject, entity)
					if err != nil {
						l.Errorf("send mail Error:%v", err)
					}
				}
			}
		}
	}()
}

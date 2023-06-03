package log

import (
	"context"
	"io"
	"time"

	"github.com/sirupsen/logrus"
)

var logger = logrus.New()

func WithField(key string, value interface{}) *logrus.Entry {
	return logger.WithField(key, value)
}

func WithFields(fields logrus.Fields) *logrus.Entry {
	return logger.WithFields(fields)
}

func WithError(err error) *logrus.Entry {
	return logger.WithError(err)
}

func WithContext(ctx context.Context) *logrus.Entry {
	return logger.WithContext(ctx)
}

func WithTime(t time.Time) *logrus.Entry {
	return logger.WithTime(t)
}

func Logf(level logrus.Level, format string, args ...interface{}) {
	logger.Logf(level, format, args...)
}

func Tracef(format string, args ...interface{}) {
	logger.Tracef(format, args...)
}

func Debugf(format string, args ...interface{}) {
	logger.Debugf(format, args...)
}

func Infof(format string, args ...interface{}) {
	logger.Infof(format, args...)
}

func Printf(format string, args ...interface{}) {
	logger.Printf(format, args...)
}

func Warnf(format string, args ...interface{}) {
	logger.Warnf(format, args...)
}

func Warningf(format string, args ...interface{}) {
	logger.Warningf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	logger.Errorf(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	logger.Fatalf(format, args...)
}

func Panicf(format string, args ...interface{}) {
	logger.Panicf(format, args...)
}

func Log(level logrus.Level, args ...interface{}) {
	logger.Log(level, args...)
}

func LogFn(level logrus.Level, fn logrus.LogFunction) {
	logger.LogFn(level, fn)
}

func Trace(args ...interface{}) {
	logger.Trace(args...)
}

func Debug(args ...interface{}) {
	logger.Debug(args...)
}

func Info(args ...interface{}) {
	logger.Info(args...)
}

func Print(args ...interface{}) {
	logger.Print(args...)
}

func Warn(args ...interface{}) {
	logger.Warn(args...)
}

func Warning(args ...interface{}) {
	logger.Warning(args...)
}

func Error(args ...interface{}) {
	logger.Error(args...)
}

func Fatal(args ...interface{}) {
	logger.Fatal(args...)
}

func Panic(args ...interface{}) {
	logger.Panic(args...)
}

func TraceFn(fn logrus.LogFunction) {
	logger.TraceFn(fn)
}

func DebugFn(fn logrus.LogFunction) {
	logger.DebugFn(fn)
}

func InfoFn(fn logrus.LogFunction) {
	logger.InfoFn(fn)
}

func PrintFn(fn logrus.LogFunction) {
	logger.PrintFn(fn)
}

func WarnFn(fn logrus.LogFunction) {
	logger.WarnFn(fn)
}

func WarningFn(fn logrus.LogFunction) {
	logger.WarningFn(fn)
}

func ErrorFn(fn logrus.LogFunction) {
	logger.ErrorFn(fn)
}

func FatalFn(fn logrus.LogFunction) {
	logger.FatalFn(fn)
}

func PanicFn(fn logrus.LogFunction) {
	logger.PanicFn(fn)
}

func Logln(level logrus.Level, args ...interface{}) {
	logger.Logln(level, args...)
}

func Traceln(args ...interface{}) {
	logger.Traceln(args...)
}

func Debugln(args ...interface{}) {
	logger.Debugln(args...)
}

func Infoln(args ...interface{}) {
	logger.Infoln(args...)
}

func Println(args ...interface{}) {
	logger.Println(args...)
}

func Warnln(args ...interface{}) {
	logger.Warnln(args...)
}

func Warningln(args ...interface{}) {
	logger.Warningln(args...)
}

func Errorln(args ...interface{}) {
	logger.Errorln(args...)
}

func Fatalln(args ...interface{}) {
	logger.Fatalln(args...)
}

func Panicln(args ...interface{}) {
	logger.Panicln(args...)
}

func Exit(code int) {
	logger.Exit(code)
}

func SetNoLock() {
	logger.SetNoLock()
}

func ParseLevel(level string) (logrus.Level, error) {
	return logrus.ParseLevel(level)
}

func SetLevel(level logrus.Level) {
	logger.SetLevel(level)
}

func GetLevel() logrus.Level {
	return logger.GetLevel()
}

func AddHook(hook logrus.Hook) {
	logger.AddHook(hook)
}

func IsLevelEnabled(level logrus.Level) bool {
	return logger.IsLevelEnabled(level)
}

func SetFormatter(formatter logrus.Formatter) {
	logger.SetFormatter(formatter)
}

func SetOutput(output io.Writer) {
	logger.SetOutput(output)
}

func SetReportCaller(reportCaller bool) {
	logger.SetReportCaller(reportCaller)
}

func ReplaceHooks(hooks logrus.LevelHooks) logrus.LevelHooks {
	return logger.ReplaceHooks(hooks)
}

func SetBufferPool(pool logrus.BufferPool) {
	logger.SetBufferPool(pool)
}

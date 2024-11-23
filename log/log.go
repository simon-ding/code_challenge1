package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

var sugar *zap.SugaredLogger

func init() {
	atom := zap.NewAtomicLevel()
	atom.SetLevel(zap.DebugLevel)

	w := zapcore.Lock(os.Stdout)

	consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())

	logger := zap.New(zapcore.NewCore(consoleEncoder, w, atom), zap.AddCallerSkip(1), zap.AddCaller())

	sugar = logger.Sugar()
}

func Logger() *zap.SugaredLogger {
	return sugar
}

func Info(args ...interface{}) {
	sugar.Info(args...)
}

func Debug(args ...interface{}) {
	sugar.Debug(args...)
}

func Warn(args ...interface{}) {
	sugar.Warn(args...)
}

func Error(args ...interface{}) {
	sugar.Error(args...)
}

func Panic(args ...interface{}) {
	sugar.Panic(args...)
}

func Infof(template string, args ...interface{}) {
	sugar.Infof(template, args...)
}

func Debugf(template string, args ...interface{}) {
	sugar.Debugf(template, args...)
}

func Warnf(template string, args ...interface{}) {
	sugar.Warnf(template, args...)
}
func Errorf(template string, args ...interface{}) {
	sugar.Errorf(template, args...)
}

func Panicf(template string, args ...interface{}) {
	sugar.Panicf(template, args...)
}

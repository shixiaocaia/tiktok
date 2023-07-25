package log

import "go.uber.org/zap"

var log *zap.Logger

func InitLog() (err error) {
	log, err = zap.NewProduction()
	return
}

func Info(msg string, fields ...zap.Field) {
	log.Info(msg, fields...)
}

func Debug(msg string, fields ...zap.Field) {
	log.Debug(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	log.Error(msg, fields...)
}

func Fatalf(msg string, fields ...zap.Field) {
	log.Fatal(msg, fields...)
}

func Sync() {
	log.Sync()
}

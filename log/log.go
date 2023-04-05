package log

import "go.uber.org/zap"

type Logger struct {
	log *zap.Logger
}

var logger *Logger = nil

func DefaultLogger() *Logger {
	if logger == nil {
		l, _ := zap.NewProduction()
		logger = &Logger{log: l}
	}

	return logger
}

func (log *Logger) Error(event string, code int, desc string) {
	log.log.Error(event,
		zap.Int("code", code),
		zap.String("desc", desc),
	)

	log.log.Sync()
}

func (log *Logger) Info(event string, code int, desc string) {
	log.log.Info(event,
		zap.Int("code", code),
		zap.String("desc", desc),
	)

	log.log.Sync()
}

func (log *Logger) Warn(event string, code int, desc string) {
	log.log.Warn(event,
		zap.Int("code", code),
		zap.String("desc", desc),
	)

	log.log.Sync()
}

func (log *Logger) Fatal(event string, code int, desc string) {
	log.log.Fatal(event,
		zap.Int("code", code),
		zap.String("desc", desc),
	)

	log.log.Sync()
}

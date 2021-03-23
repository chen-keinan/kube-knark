package logger

import "go.uber.org/zap"

//NewZapLogger zap logger object
func NewZapLogger() *zap.Logger {
	logger, err := zap.NewProduction()
	if err != nil {
		panic("failed to create zap logger instance")
	}
	return logger
}

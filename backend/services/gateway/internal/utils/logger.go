package utils

import "go.uber.org/zap"

var logger *zap.Logger

func SetLogger(l *zap.Logger) {
	logger = l
}

func Logger() *zap.Logger {
	return logger
}

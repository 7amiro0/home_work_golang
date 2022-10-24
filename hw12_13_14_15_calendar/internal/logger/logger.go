package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	logger *zap.SugaredLogger
}

func configLogger(level string) zap.Config {
	atomicLevel, err := zap.ParseAtomicLevel(level)
	if err != nil {
		panic(err)
	}

	config := zap.Config{
		Encoding:    "console",
		Level:       atomicLevel,
		OutputPaths: []string{"stdout"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey: "message",
		},
	}

	return config
}

func New(level string) *Logger {
	config := configLogger(level)

	logger, err := config.Build()
	if err != nil {
		panic(err)
	}

	return &Logger{
		logger: logger.Sugar(),
	}
}

func (l Logger) Fatal(msg ...any) {
	l.logger.Fatal(msg...)
}

func (l Logger) Info(msg ...any) {
	l.logger.Info(msg...)
}

func (l Logger) Error(msg ...any) {
	l.logger.Error(msg...)
}

func (l Logger) Debug(msg ...any) {
	l.logger.Debug(msg...)
}

func (l Logger) Warn(msg ...any) {
	l.logger.Warn(msg...)
}

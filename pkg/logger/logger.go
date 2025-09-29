package logger

import (
	"log"
	"os"
	"rtk/api-mocker/internal/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger interface {
	Debugw(msg string, keysAndValues ...any)
	Infow(msg string, keysAndValues ...any)
	Info(args ...any)
	Warnw(msg string, keysAndValues ...any)
	Errorw(msg string, keysAndValues ...any)
	Errorf(template string, args ...any)
	Infof(template string, args ...any)
}

func New(config *config.Config) *zap.SugaredLogger {
	level := zap.NewAtomicLevel()

	level.SetLevel(zap.DebugLevel)

	consoleEncoderCfg := zap.NewDevelopmentEncoderConfig()

	consoleEncoderCfg.StacktraceKey = ""
	consoleEncoderCfg.CallerKey = "C"
	consoleEncoderCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	consoleEncoderCfg.EncodeTime = zapcore.RFC3339TimeEncoder

	consoleEncoder := zapcore.NewConsoleEncoder(consoleEncoderCfg)

	// file encoder (JSON)
	fileEncoderCfg := zap.NewProductionEncoderConfig()
	fileEncoderCfg.EncodeTime = zapcore.RFC3339TimeEncoder
	fileEncoder := zapcore.NewJSONEncoder(fileEncoderCfg)

	file, err := os.OpenFile(config.LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		log.Fatalf("failed to open log file: %v", err)
	}

	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), level),
		zapcore.NewCore(fileEncoder, zapcore.AddSync(file), level),
	)

	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))
	return logger.Sugar()
}

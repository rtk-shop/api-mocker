package logger

import (
	"log"
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
}

func New(config *config.Config) *zap.SugaredLogger {
	level := zap.NewAtomicLevel()

	// if config.isDev {
	//  level.SetLevel(zap.DebugLevel)
	// } else {
	// 	level.SetLevel(zap.InfoLevel)
	// }

	level.SetLevel(zap.DebugLevel)

	cfg := zap.Config{
		Level:            level,
		Development:      true,
		Encoding:         "console",
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: func() zapcore.EncoderConfig {
			encCfg := zap.NewDevelopmentEncoderConfig()

			encCfg.StacktraceKey = ""
			encCfg.CallerKey = "C"
			encCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
			encCfg.EncodeTime = zapcore.RFC3339TimeEncoder

			return encCfg
		}(),
	}

	logger, err := cfg.Build()
	if err != nil {
		log.Fatalf("failed to initialize logger: %v", err)
	}

	return logger.Sugar()
}

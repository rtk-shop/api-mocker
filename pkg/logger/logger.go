package logger

import (
	"rtk/api-mocker/internal/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New(config *config.Config) (*zap.SugaredLogger, error) {
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
		return nil, err
	}

	return logger.Sugar(), nil
}

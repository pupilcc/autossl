package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func RequestLogger() echo.MiddlewareFunc {
	logger := GetLogger()

	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogMethod: true,
		LogStatus: true,
		LogError:  true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			time := zap.String("time", v.StartTime.Format("2006-01-02 15:04:05"))
			uri := zap.String("URI", v.URI)
			method := zap.String("method", v.Method)
			status := zap.Int("status", v.Status)
			if v.Error != nil {
				logger.Error("request",
					time,
					uri,
					method,
					status,
					zap.String("error", v.Error.Error()),
				)
			} else {
				logger.Info("request",
					time,
					uri,
					method,
					status,
				)
			}
			return nil
		},
	})
}

func GetLogger() *zap.Logger {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.EpochTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(zap.InfoLevel)

	config := zap.Config{
		Level:            atomicLevel,
		Development:      false,
		Sampling:         nil,
		Encoding:         "console", // "json" or "console"
		EncoderConfig:    encoderConfig,
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}
	config.DisableStacktrace = true
	logger, _ := config.Build()
	return logger
}

package logger

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// The log format can either be text or JSON.
const (
	JSONFormat = "json"
	TextFormat = "text"
)

const (
	RFC3339Z = "2006-01-02T15:04:05.000Z" // time.RFC3339 with Zulu time zone
)

// Config stores the config for the logger.
type Config struct {
	Level  string `json:"level" envconfig:"LOG_LEVEL" default:"info"`
	Format string `json:"format" envconfig:"LOG_FORMAT" default:"json"`
}

func toZapLevel(level string) (zapcore.Level, error) {
	var l zapcore.Level
	err := l.UnmarshalText([]byte(level))
	return l, err
}

func getEncoder(logFormat string) zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	// ISO8601-formatted string in UTC/Zulu timezone
	encoderConfig.EncodeTime = zapcore.TimeEncoder(func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.UTC().Format(RFC3339Z))
	})
	encoderConfig.TimeKey = "time" // This will change the key from 'ts' to 'time'

	if JSONFormat == logFormat {
		return zapcore.NewJSONEncoder(encoderConfig)
	}

	return zapcore.NewConsoleEncoder(encoderConfig)
}

func New(config *Config) (*zap.SugaredLogger, error) {
	cores := []zapcore.Core{}

	level, err := toZapLevel(config.Level)
	if err != nil {
		return nil, err
	}
	encoder := getEncoder(config.Format)
	writer := zapcore.Lock(os.Stdout)
	core := zapcore.NewCore(encoder, writer, level)
	cores = append(cores, core)

	combinedCore := zapcore.NewTee(cores...)

	zLogger := zap.New(combinedCore,
		zap.AddCaller(),
	).Sugar()

	return zLogger, nil
}

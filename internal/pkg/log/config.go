package log

import "go.uber.org/zap/zapcore"

type (
	StringEncode string

	Config struct {
		Encoding string       `env:"LOG_ENCODING" envDefault:"console"`
		Level    StringEncode `env:"LOG_LEVEL" envDefault:"debug"`
		Enabled  bool         `env:"LOG_ENABLED" envDefault:"true"`
	}
)

func (s StringEncode) Level() zapcore.Level {
	switch s {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		panic("Unmapped log level, received: " + s)
	}
}

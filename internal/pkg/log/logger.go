package log

import (
	stdLog "log"

	"github.com/caarlos0/env/v9"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	writer       *zap.Logger
	loggerConfig zap.Config
)

const (
	DebugLevel  = Level(zap.DebugLevel)
	InfoLevel   = Level(zap.InfoLevel)
	WarnLevel   = Level(zap.WarnLevel)
	ErrorLevel  = Level(zap.ErrorLevel)
	DPanicLevel = Level(zap.DPanicLevel)
	PanicLevel  = Level(zap.PanicLevel)
	FatalLevel  = Level(zap.FatalLevel)
	MinLevel    = DebugLevel
	MaxLevel    = FatalLevel
)

var (
	Any      = zap.Any
	Binary   = zap.Binary
	Bool     = zap.Bool
	Float32  = zap.Float32
	Float64  = zap.Float64
	Int      = zap.Int
	Int8     = zap.Int8
	Int16    = zap.Int16
	Int32    = zap.Int32
	Int64    = zap.Int64
	String   = zap.String
	Stringer = zap.Stringer
	Uint     = zap.Uint
	Uint8    = zap.Uint8
	Uint16   = zap.Uint16
	Uint32   = zap.Uint32
	Uint64   = zap.Uint64
	StdError = zap.Error
	Time     = zap.Time
	Duration = zap.Duration
	Object   = zap.Object
	Strings  = zap.Strings
)

type (
	Level  int8
	Logger interface {
		Error(msg string, field ...Field)
		Info(msg string, field ...Field)
		Warn(msg string, field ...Field)
		Debug(msg string, field ...Field)
		Fatal(msg string, field ...Field)
		Panic(msg string, field ...Field)
		With(fields ...Field) *zap.Logger
	}
	Field = zap.Field
)

func init() {
	conf := new(Config)
	err := env.Parse(conf)
	if err != nil {
		panic(err.Error())
	}

	zap.NewProductionConfig()

	// check if logging is enabled
	var outputPaths, errorOutputPaths []string
	if conf.Enabled {
		outputPaths = append(outputPaths, "stdout")
		errorOutputPaths = append(errorOutputPaths, "stderr")
	}

	loggerConfig = zap.Config{
		Level:            zap.NewAtomicLevelAt(conf.Level.Level()),
		Encoding:         conf.Encoding,
		OutputPaths:      outputPaths,
		ErrorOutputPaths: errorOutputPaths,
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "msg",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
		},
		DisableStacktrace: true,
	}

	l, err := loggerConfig.Build(zap.AddCallerSkip(1))
	if err != nil {
		panic(err)
	}

	writer = l
	zap.ReplaceGlobals(l)
}

func Error(msg string, field ...Field) {
	zap.L().Error(msg, field...)
}

func Info(msg string, field ...Field) {
	zap.L().Info(msg, field...)
}

func Warn(msg string, field ...Field) {
	zap.L().Warn(msg, field...)
}

func Debug(msg string, field ...Field) {
	zap.L().Debug(msg, field...)
}

func Fatal(msg string, field ...Field) {
	zap.L().Fatal(msg, field...)
}

func Panic(msg string, field ...Field) {
	zap.L().Panic(msg, field...)
}

type logWrapper struct {
	*zap.Logger
}

func With(fields ...Field) Logger {
	return logWrapper{zap.L().With(fields...)}
}

func NewLogger() Logger {
	return writer
}

func Base() *zap.Logger {
	return writer
}

func SetLogLevel(l Level) {
	loggerConfig.Level.SetLevel(zapcore.Level(l))
}

func NewStdLogger() *stdLog.Logger {
	return zap.NewStdLog(zap.L())
}

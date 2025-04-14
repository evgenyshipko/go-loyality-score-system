package logger

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

var Instance *zap.SugaredLogger

const (
	colorCyan = "\033[36m"
)

func colorTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	coloredTime := fmt.Sprintf("%s%s", colorCyan, t.Format("2006-01-02 15:04:05"))
	enc.AppendString(coloredTime)
}

var logger *zap.Logger

// ЗАПОМНИТЬ! init вызывается при старте пакета, инициализируя логгер
func init() {
	encoderCfg := zapcore.EncoderConfig{
		TimeKey:          "time",
		LevelKey:         "level",
		NameKey:          "logger",
		CallerKey:        "caller",
		MessageKey:       "msg",
		StacktraceKey:    "stacktrace",
		EncodeTime:       colorTimeEncoder,
		EncodeLevel:      zapcore.CapitalColorLevelEncoder,
		EncodeCaller:     zapcore.ShortCallerEncoder,
		ConsoleSeparator: " | ",
	}

	consoleEncoder := zapcore.NewConsoleEncoder(encoderCfg)

	core := zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), zapcore.DebugLevel)

	logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	Instance = logger.Sugar()
}

func Sync() {
	if logger != nil {
		_ = logger.Sync()
	}
}

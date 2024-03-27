package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"time"
)

var Logger *zap.Logger

func Init() {
	Logger = NewLogger("default")
	zap.ReplaceGlobals(Logger)
}

func NewLogger(field string) *zap.Logger {
	lumberjackLogger := &lumberjack.Logger{ //日志切割
		Filename:   "/var/log/unionAuth/ua.log",
		MaxSize:    10,
		MaxBackups: 3,
		MaxAge:     28,
		Compress:   true,
	}

	consoleEncoder := zap.NewProductionEncoderConfig()
	consoleEncoder.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05")) //格式化时间
	}
	consoleEncoder.EncodeLevel = zapcore.CapitalColorLevelEncoder //让level字段变成彩色
	consoleEncoder.EncodeCaller = func(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(field)
		enc.AppendString(caller.String())
	}

	fileEncoder := zap.NewProductionEncoderConfig()
	fileEncoder.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05"))
	}
	fileEncoder.EncodeLevel = zapcore.CapitalLevelEncoder
	fileEncoder.EncodeCaller = func(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(field)
		enc.AppendString(caller.String())
	}

	core := zapcore.NewTee(
		zapcore.NewCore(zapcore.NewJSONEncoder(fileEncoder), zapcore.AddSync(lumberjackLogger), zap.InfoLevel),
		zapcore.NewCore(zapcore.NewConsoleEncoder(consoleEncoder), zapcore.Lock(os.Stdout), zap.DebugLevel),
	)

	Logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))

	defer Logger.Sync()
	return Logger
}

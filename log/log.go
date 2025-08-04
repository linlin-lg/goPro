package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

var Logger *zap.Logger

func InitLog() {
	core := zapcore.NewCore(getEncoder(), getLogWriter(), zapcore.DebugLevel)
	Logger = zap.New(core)

	//sugarLogger := Logger.Sugar()
	//l := 10
	//sugarLogger.Debugf("d%123", l)
}

func getLogWriter() zapcore.WriteSyncer {
	hook := lumberjack.Logger{
		Filename: "./log/pro.log",
		MaxSize: 128,
		MaxBackups: 30,
		MaxAge: 30,
		Compress: true,
	}
	writeSyncer := zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&hook))
	return writeSyncer
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey: "time",
		LevelKey: "lever",
		NameKey: "logger",
		CallerKey: "caller",
		StacktraceKey: "stacktraceKey",
		LineEnding: zapcore.DefaultLineEnding,
		EncodeLevel: zapcore.LowercaseLevelEncoder,
		EncodeTime: zapcore.RFC3339TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller: zapcore.FullCallerEncoder,
	}
	return zapcore.NewJSONEncoder(encoderConfig)
}
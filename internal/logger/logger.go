package logger

import (
	"os"
	"rd-read-book-project/global"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func InitLogger() {
	// 1️⃣ 编码器配置
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoder := zapcore.NewJSONEncoder(encoderConfig)

	// 2️⃣ info 日志文件
	infoWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "./logs/info.log",
		MaxSize:    10,
		MaxBackups: 15,
		MaxAge:     90,
		Compress:   true,
	})

	// 3️⃣ error 日志文件
	errorWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "./logs/error.log",
		MaxSize:    10,
		MaxBackups: 15,
		MaxAge:     90,
		Compress:   true,
	})

	// 4️⃣ 控制等级过滤

	// info 级别：Debug ~ Info ~ Warn
	infoLevel := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level < zapcore.ErrorLevel
	})

	// error 级别：Error 及以上
	errorLevel := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level >= zapcore.ErrorLevel
	})

	// 5️⃣ 创建两个 core
	infoCore := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(
		zapcore.AddSync(os.Stdout),
		infoWriter,
	), infoLevel)

	errorCore := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(
		zapcore.AddSync(os.Stdout),
		errorWriter,
	), errorLevel)

	// 6️⃣ 合并 core
	core := zapcore.NewTee(infoCore, errorCore)

	// 7️⃣ 创建 logger
	logger := zap.New(core, zap.AddCaller())

	global.ZapLogger = logger
}

func Debug(msg string, fields ...zap.Field) {
	global.ZapLogger.Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	global.ZapLogger.Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	global.ZapLogger.Warn(msg, fields...)
}

func Error(msg string, err error, fields ...zap.Field) {
	if err != nil {
		fields = append(fields, zap.Error(err))
	}
	global.ZapLogger.Error(msg, fields...)
}

func Fatal(msg string, err error, fields ...zap.Field) {
	if err != nil {
		fields = append(fields, zap.Error(err))
	}
	global.ZapLogger.Fatal(msg, fields...)
}

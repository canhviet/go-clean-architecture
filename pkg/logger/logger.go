package logger

import (
    "github.com/spf13/viper"
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
)

var Log *zap.Logger

func Init() {
    level := viper.GetString("log.level")
    var zapLevel zapcore.Level
    zapLevel.UnmarshalText([]byte(level))

    config := zap.NewProductionConfig()
    config.Level = zap.NewAtomicLevelAt(zapLevel)
    config.EncoderConfig.TimeKey = "timestamp"
    config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

    var err error
    Log, err = config.Build()
    if err != nil {
        panic(err)
    }
}

func Info(msg string, fields ...zap.Field)  { Log.Info(msg, fields...) }
func Error(msg string, fields ...zap.Field) { Log.Error(msg, fields...) }
func Fatal(msg string, fields ...zap.Field) { Log.Fatal(msg, fields...) }
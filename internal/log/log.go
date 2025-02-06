// internal/log/log.go
package log

import (
	"aidan/phone/pkg/util"
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var logger *zap.Logger

func init() {
	var err error
	workingDir, err := util.GetWorkingDir()
	if err != nil {
		panic(fmt.Errorf("unable to get working dir: %w", err))
	}

	stdout := zapcore.AddSync(os.Stdout)

	//auto log rotate
	file := zapcore.AddSync(&lumberjack.Logger{
		Filename:   workingDir + "/" + "phone.log",
		MaxSize:    10, // megabytes
		MaxBackups: 3,
		MaxAge:     30, // days
	})

	level := zap.NewAtomicLevelAt(zap.InfoLevel)

	productionCfg := zap.NewProductionEncoderConfig()
	productionCfg.TimeKey = "timestamp"
	productionCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	developmentCfg := zap.NewDevelopmentEncoderConfig()
	developmentCfg.EncodeLevel = zapcore.CapitalLevelEncoder

	consoleEncoder := zapcore.NewConsoleEncoder(developmentCfg)
	fileEncoder := zapcore.NewJSONEncoder(productionCfg)

	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, stdout, level),
		zapcore.NewCore(fileEncoder, file, level),
	)

	l := zap.New(core)
	logger = l
}

func GetLogger() *zap.Logger {
	return logger
}

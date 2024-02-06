package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func createLogger(logDirPath string, fileName string) *zap.Logger {
	stdout := zapcore.AddSync(os.Stdout)

	//auto log rotate
	file := zapcore.AddSync(&lumberjack.Logger{
		Filename:   logDirPath + "/" + fileName,
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

	return zap.New(core)
}

func getWorkingDir() (string, error) {
	return filepath.Abs(filepath.Dir(os.Args[0]))
}

// json formatted config loader, pass a config pointer
func loadConfig(path string, config interface{}) error {
	configFile, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open json config file: %w", err)
	}
	defer configFile.Close()

	data, err := ioutil.ReadAll(configFile)
	if err != nil {
		return fmt.Errorf("failed to read json config file: %w", err)
	}

	err = json.Unmarshal(data, config)
	if err != nil {
		return fmt.Errorf("failed to unmarshal json config file: %w", err)
	}

	return nil
}

func getHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		return "Unknown"
	}
	return hostname
}

func getStackTrace() string {
	// Retrieve the stack trace
	stack := make([]byte, 4096)
	length := runtime.Stack(stack, false)

	// Convert the stack trace to a string
	return string(stack[:length])
}

// depth of 1 is the file that called this, 1 < n = files higher in the stack
func getCaller(depth int) string {
	pc, file, line, ok := runtime.Caller(depth)

	if !ok {
		file = "unknown"
		line = 0
	} else {
		file = file[strings.LastIndex(file, "/")+1:]
	}

	funcName := runtime.FuncForPC(pc).Name()

	return fmt.Sprintf("[%s] %s:%d\n", funcName, file, line)
}

func removeElementFromSlice(slicePtr interface{}, indexToRemove int) error {
	sliceValue := reflect.ValueOf(slicePtr).Elem()

	if sliceValue.Kind() != reflect.Slice {
		return errors.New("cannot remove element from invalid slice")
	}

	if indexToRemove < 0 || indexToRemove >= sliceValue.Len() {
		return errors.New("index of element to remove is out of range")
	}

	// Create a new slice by excluding the element at the specified index
	newSlice := reflect.AppendSlice(sliceValue.Slice(0, indexToRemove), sliceValue.Slice(indexToRemove+1, sliceValue.Len()))

	sliceValue.Set(newSlice)

	return nil
}

func localTimeNow() (time.Time, error) {
	location, err := time.LoadLocation("America/Chicago")
	if err != nil {
		return time.Time{}, err
	}

	currentTime := time.Now().In(location)
	return currentTime, nil
}

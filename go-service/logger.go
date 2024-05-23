package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
)

var logger = NewMultiLogger()

type MultiLogger struct {
	consoleLogger *logrus.Logger
	fileLogger    *logrus.Logger
}

func NewMultiLogger() *MultiLogger {
	// Init loggers
	logger := &MultiLogger{
		consoleLogger: logrus.New(),
		fileLogger:    logrus.New(),
	}
	// Set log Output Stream
	logger.consoleLogger.SetOutput(os.Stdout)
	//logFile, err := os.OpenFile(filepath.Join(filePath), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	//if err != nil {
	//	logger.consoleLogger.Fatalln("Failed to open log file:", err)
	//}
	logger.fileLogger.SetOutput(io.Discard)
	// Set default log level
	logger.consoleLogger.SetLevel(logrus.InfoLevel)
	logger.fileLogger.SetLevel(logrus.DebugLevel)
	return logger
}

func (logger *MultiLogger) UpdateFile(filePath string) {
	logFile, err := os.OpenFile(filepath.Join(filePath), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logger.consoleLogger.Fatalln("Failed to open log file:", err)
	}
	logger.fileLogger.SetOutput(logFile)
}

func (logger *MultiLogger) Trace(args ...interface{}) {
	logger.consoleLogger.Trace(args...)
	logger.fileLogger.Trace(args...)
}

func (logger *MultiLogger) Debug(args ...interface{}) {
	logger.consoleLogger.Debug(args...)
	logger.fileLogger.Debug(args...)
}

func (logger *MultiLogger) Info(args ...interface{}) {
	logger.consoleLogger.Info(args...)
	logger.fileLogger.Info(args...)
}

func (logger *MultiLogger) Warn(args ...interface{}) {
	logger.consoleLogger.Warn(args...)
	logger.fileLogger.Warn(args...)
}

func (logger *MultiLogger) Error(args ...interface{}) {
	logger.consoleLogger.Error(args...)
	logger.fileLogger.Error(args...)
}

func (logger *MultiLogger) Fatal(args ...interface{}) {
	logger.consoleLogger.Fatal(args...)
	logger.fileLogger.Fatal(args...)
}

func (logger *MultiLogger) Panic(args ...interface{}) {
	logger.consoleLogger.Panic(args...)
	logger.fileLogger.Panic(args...)
}

func (logger *MultiLogger) SetLevel(level string) {
	if level == "default" {
		return
	}
	logLevel, err := GetLogLevelFromName(level)
	if err != nil {
		logger.Warn(err)
	} else {
		logger.consoleLogger.SetLevel(logLevel)
		logger.fileLogger.SetLevel(logLevel)
	}
}

func GetLogLevelFromName(level string) (logrus.Level, error) {
	var logLevel logrus.Level
	switch level {
	case "trace":
		logLevel = logrus.TraceLevel
	case "debug":
		logLevel = logrus.DebugLevel
	case "info":
		logLevel = logrus.InfoLevel
	case "warning":
		logLevel = logrus.WarnLevel
	case "error":
		logLevel = logrus.ErrorLevel
	case "fatal":
		logLevel = logrus.FatalLevel
	case "panic":
		logLevel = logrus.PanicLevel
	default:
		return logrus.InfoLevel, fmt.Errorf("unknown log level: %s", level)
	}
	return logLevel, nil
}

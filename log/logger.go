package log

import (
	"fmt"
	"time"
)

type Logger struct {
	filePath string
	fileName string
	fileCount int
	singleFileSize int
}

func NewLogger() *Logger {
	return &Logger{}
}

func (logger *Logger)filePath(path string) *Logger {
	logger.filePath = path
	return logger
}

func (logger *Logger)fileName(name string) *Logger {
	logger.fileName = name
	return logger
}

func (logger *Logger)fileCount(count int) *Logger {
	logger.fileCount = count
	return logger
}

func (logger *Logger)SingleFileSize(size int) *Logger {
	logger.singleFileSize = size
	return logger
}

func log(level, tag, msg string) {
	fmt.Printf("%s %s/%s: %s\n", time.Now().Format("2006-01-02 15:04:05"), level, tag, msg)
}

func (logger *Logger)D(tag, format string, args ...interface{})  {
	log("D", tag, fmt.Sprintf(format, args...))
}

func (logger *Logger)I(tag, format string, args ...interface{})  {
	log("I", tag, fmt.Sprintf(format, args...))
}

func (logger *Logger)W(tag, format string, args ...interface{})  {
	log("W", tag, fmt.Sprintf(format, args...))
}

func (logger *Logger)E(tag, format string, args ...interface{})  {
	log("E", tag, fmt.Sprintf(format, args...))
}



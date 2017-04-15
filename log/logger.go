package log

import (
	"fmt"
	"time"
)

type Logger struct {
	filePath       string
	fileName       string
	fileCount      int
	singleFileSize int64
	cacheCount     int
	cache          chan string
	isWriting      bool
}

func NewLogger() *Logger {
	return &Logger{}
}

func (logger *Logger) FilePath(path string) *Logger {
	logger.filePath = path
	return logger
}

func (logger *Logger) FileName(name string) *Logger {
	logger.fileName = name
	return logger
}

func (logger *Logger) FileCount(count int) *Logger {
	logger.fileCount = count
	return logger
}

func (logger *Logger) SingleFileSize(size int64) *Logger {
	logger.singleFileSize = size
	return logger
}

func (logger *Logger) CacheCount(count int) *Logger {
	logger.cacheCount = count
	return logger
}

func (logger *Logger) log(level, tag, msg string) {
	log := fmt.Sprintf("%s %s/%s: %s\n", time.Now().Format("2006-01-02 15:04:05"), level, tag, msg)
	fmt.Print(log)
	if logger.fileName != "" {
		logger.cache <- log
		if len(logger.cache) >= logger.cacheCount && !logger.isWriting {
			logger.isWriting = true
			go func() {
				if logger.filePath == "" {
					logger.filePath = "./"
				}
				//TODO
			}()
		}
	}
}

func (logger *Logger) D(tag, format string, args ...interface{}) {
	logger.log("D", tag, fmt.Sprintf(format, args...))
}

func (logger *Logger) I(tag, format string, args ...interface{}) {
	logger.log("I", tag, fmt.Sprintf(format, args...))
}

func (logger *Logger) W(tag, format string, args ...interface{}) {
	logger.log("W", tag, fmt.Sprintf(format, args...))
}

func (logger *Logger) E(tag, format string, args ...interface{}) {
	logger.log("E", tag, fmt.Sprintf(format, args...))
}

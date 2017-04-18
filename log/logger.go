package log

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type Logger struct {
	path           string
	singleFileSize int64
	cache          chan string
	isWriting      bool
	exit           bool
	mutex          sync.Mutex
}

type LogBuilder struct {
	path           string
	singleFileSize int64
}

func NewLogBuilder() *LogBuilder {
	return &LogBuilder{path: "./", singleFileSize: 4096}
}

func (builder *LogBuilder) FilePath(path string) *LogBuilder {
	builder.path = path
	return builder
}

func (builder *LogBuilder) SingleFileSize(size int64) *LogBuilder {
	builder.singleFileSize = size
	return builder
}

func (builder *LogBuilder) build() *Logger {
	return &Logger{path: builder.path,
		singleFileSize: builder.singleFileSize,
		cache:          make(chan string, 100)}
}

func zipFile(zipFileName string, pFile *os.File) error {
	zipFile, err := os.Create(zipFileName)
	if err != nil {
		return err
	}
	zw := zip.NewWriter(zipFile)
	defer zw.Close()
	fw, err := zw.Create(pFile.Name())
	io.Copy(fw, pFile)
	fmt.Println("zipFile err:", err)
	return err
}

func findLastModified(dir string) os.FileInfo {
	var fileInfo os.FileInfo
	fileInfoSlice, err := ioutil.ReadDir(dir)
	if err == nil {
		for _, item := range fileInfoSlice {
			if item.IsDir() || !strings.HasSuffix(item.Name(), ".log") {
				continue
			}
			if fileInfo == nil {
				fileInfo = item
			} else if item.ModTime().Unix() > fileInfo.ModTime().Unix() {
				fileInfo = item
			}

		}
		if fileInfo == nil {
			pFile, err := os.Create(dir + string(filepath.Separator) + time.Now().Format("2006-01-02 15:04:05") + ".log")
			if err == nil {
				fileInfo, err = os.Stat(pFile.Name())
			}
		}
	}
	return fileInfo
}

func (logger *Logger) log(level, tag, msg string) {
	log := fmt.Sprintf("%s %s/%s: %s\n", time.Now().Format("2006-01-02 15:04:05"), level, tag, msg)
	fmt.Print(log)
	if logger.exit {
		return
	}
	logger.cache <- log
	if logger.isWriting {
		return
	}
	logger.mutex.Lock()
	defer logger.mutex.Unlock()
	if !logger.isWriting {
		logger.isWriting = true
		go func() {
			var pFile *os.File
			for log := range logger.cache {
				if pFile == nil {
					fileInfo := findLastModified(logger.path)
					file, err := os.OpenFile(fileInfo.Name(), os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
					if err != nil {
						continue
					}
					pFile = file
				}
				pFile.WriteString(log)
				fileInfo, _ := os.Stat(pFile.Name())
				if fileInfo != nil && fileInfo.Size() > logger.singleFileSize {
					idx := strings.LastIndex(fileInfo.Name(), string(filepath.Separator))
					var dir string
					if idx == -1 {
						dir = "./"
					} else {
						dir = fileInfo.Name()[:idx]
					}
					zipFileName := dir + "log-" + time.Now().Format("2006-01-02 15:04:05") + ".zip"
					pFile.Seek(0, os.SEEK_SET)
					zipFile(zipFileName, pFile)
					os.Remove(pFile.Name())
					pFile.Close()
					pFile = nil
				}
			}
		}()
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

func (logger *Logger) Close() {
	logger.exit = true
	close(logger.cache)
}

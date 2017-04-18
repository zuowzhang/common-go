package log

import (
	"os"
	"testing"
	"time"
)

func TestZipFile(t *testing.T) {
	pFile, _ := os.Open("./logger_test.go")
	err := zipFile("/home/zbin/a.zip", pFile)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("ok")
}

func TestLogger_SingleFileSize(t *testing.T) {
	logger := NewLogBuilder().Build()
	for idx := 0; idx < 100; idx++ {
		time.Sleep(100 * time.Millisecond)
		logger.D("logger_test", "测试测试测试测试测试测试")
		if idx == 80 {
			logger.Close()
			time.Sleep(time.Second)
		}
	}
}

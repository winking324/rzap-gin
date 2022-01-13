package rzap_gin

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/winking324/rzap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

type line struct {
	Status int    `json:"status"`
	Level  string `json:"level"`
	TS     string `json:"ts"`
	Path   string `json:"path"`
	Method string `json:"method"`
	Query  string `json:"query"`
	Errors string `json:"errors"`
}

func expectRemoveFile(filePath string, t *testing.T) {
	if _, err := os.Stat(filePath); err == nil {
		if err := os.Remove(filePath); err != nil {
			t.Errorf("Remove file: %s failed", filePath)
		}
	}
}

func expectExistFile(filePath string, t *testing.T) {
	if _, err := os.Stat(filePath); err != nil {
		t.Errorf("File: %s not exist", filePath)
	}
}

func expectReadFile(filePath string, t *testing.T) *line {
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		t.Errorf("Read file: %s failed", filePath)
		return nil
	}

	l := &line{}
	if err := json.Unmarshal(b, l); err != nil {
		t.Errorf("Unmarshal line: %s failed", string(b))
		return nil
	}
	return l
}

func TestRZapGin(t *testing.T) {
	logFile := "/tmp/rzap.log"
	expectRemoveFile(logFile, t)

	rzap.NewGlobalLogger([]zapcore.Core{
		rzap.NewCore(&lumberjack.Logger{
			Filename: logFile,
		}, zap.InfoLevel),
	})

	r := gin.New()
	r.Use(Logger(nil), Recovery(nil, false))

	r.GET("/test", func(context *gin.Context) {
		context.JSON(http.StatusNoContent, nil)
	})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(res, req)

	if res.Code != http.StatusNoContent {
		t.Errorf("Response code: %d error", res.Code)
	}

	expectExistFile(logFile, t)
	l := expectReadFile(logFile, t)
	if l.Status != http.StatusNoContent {
		t.Error("Log 'status' error")
	}
	if l.Level != "INFO" {
		t.Error("Log 'level' error")
	}
	if l.Path != "/test" {
		t.Error("Log 'path' error")
	}
	if l.Method != http.MethodGet {
		t.Error("Log 'method' error")
	}
}

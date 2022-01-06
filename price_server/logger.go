package main

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"strings"
)

type Configuration struct {
	EnableConsole     bool
	ConsoleJSONFormat bool
	ConsoleLevel      string
	EnableFile        bool
	FileJSONFormat    bool
	FileLevel         string
	FileLocation      string
}

func DefaultConfiguration() *Configuration {
	return &Configuration{
		EnableConsole:     true,
		ConsoleJSONFormat: false,
		ConsoleLevel:      logrus.DebugLevel.String(),
		EnableFile:        true,
		FileJSONFormat:    false,
		FileLevel:         logrus.InfoLevel.String(),
		FileLocation:      "logs/warehouse.log",
	}
}

type MyFormatter struct{}

var levelList = []string{
	"PANIC",
	"FATAL",
	"ERROR",
	"WARN",
	"INFO",
	"DEBUG",
	"TRACE",
}

func (mf *MyFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	dataB := &bytes.Buffer{}

	for k := range entry.Data {
		v := entry.Data[k]
		dataB.WriteString(k)
		dataB.WriteByte('=')
		dataB.WriteString(fmt.Sprint(v))
		dataB.WriteByte(' ')
	}

	level := levelList[int(entry.Level)]
	strList := strings.Split(entry.Caller.File, "/")
	fileName := strList[len(strList)-1]
	if len(entry.Data) > 0 {
		b.WriteString(fmt.Sprintf("%s - [%s:%d] - %s - %s - %s\n",
			entry.Time.Format("2006-01-02 15:04:05"), fileName,
			entry.Caller.Line, level, entry.Message, dataB.String()))
	} else {
		b.WriteString(fmt.Sprintf("%s - [%s:%d] - %s - %s\n",
			entry.Time.Format("2006-01-02 15:04:05"), fileName,
			entry.Caller.Line, level, entry.Message))
	}
	return b.Bytes(), nil
}

func InitLogrusLogger(config *Configuration) error {
	logLevel := config.ConsoleLevel
	if logLevel == "" {
		logLevel = config.FileLevel
	}
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		return err
	}
	stdOutHandler := os.Stdout
	fileHandler := &lumberjack.Logger{
		Filename: config.FileLocation,
		MaxSize:  20,
		Compress: true,
		MaxAge:   30,
	}
	logrus.SetFormatter(&MyFormatter{})
	logrus.SetLevel(level)
	logrus.SetReportCaller(true)

	if config.EnableConsole && config.EnableFile {
		logrus.SetOutput(io.MultiWriter(stdOutHandler, fileHandler))
	} else if config.EnableFile {
		logrus.SetOutput(fileHandler)
	}

	return nil
}

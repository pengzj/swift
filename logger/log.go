package logger

import (
	"log"
	"time"
	"fmt"
	"sync"
	"path/filepath"
	"os"
)

type Logger struct {
	fileName string
	level int
	mu sync.RWMutex
	fp *os.File
	lastEditDay string
}

const (
	LevelTrace = iota
	LevelDebug
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

var level = LevelTrace


func SetLevel(l int)  {
	level = l
}

func GetLevel() int {
	return level
}

func (logger *Logger) Trace(v ...interface{})  {
	if level <= LevelTrace {
		write(logger, LevelTrace, v...)
	}
}

func (logger *Logger) Debug(v ...interface{})  {
	if level <= LevelDebug {
		write(logger, LevelDebug, v...)
	}
}

func (logger *Logger) Info(v ...interface{})  {
	if level <= LevelInfo {
		write(logger,LevelInfo, v...)
	}
}

func (logger *Logger) Warn(v ...interface{})  {
	if level <= LevelWarn {
		write(logger,LevelWarn, v...)
	}
}

func (logger *Logger) Error(v ...interface{})  {
	if level <= LevelError {
		write(logger,LevelError, v...)
	}
}

func (logger *Logger) Fatal(v ...interface{})  {
	if level <= LevelFatal {
		write(logger,LevelFatal, v...)
	}
}

func (logger *Logger) SetFile(fileName string)  {
	logger.mu.Lock()
	defer logger.mu.Unlock()
	logger.fileName = fileName

	dir, file := filepath.Split(logger.fileName)
	ext := filepath.Ext(file)
	var name = file[0:len(file)-len(ext)]
	now := time.Now()
	year, month, day := now.Date()
	ymd := fmt.Sprintf("%04d%02d%02d", year, month, day)

	fileName =   name + "." + ymd + ext
	fileName = filepath.Join(dir, fileName)

	fp, err := os.OpenFile(fileName, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0755)
	if err != nil {
		log.Fatal(err)
	}
	logger.fp = fp
	logger.lastEditDay = ymd
	log.SetOutput(fp)
}


func write(logger *Logger, l int, v ...interface{})  {
	logger.mu.Lock()
	defer logger.mu.Unlock()
	dir, file := filepath.Split(logger.fileName)
	ext := filepath.Ext(file)
	var name = file[0:len(file)-len(ext)]
	now := time.Now()
	year, month, day := now.Date()
	ymd := fmt.Sprintf("%04d%02d%02d", year, month, day)


	if ymd != logger.lastEditDay {
		var fileName =   name + ymd + ext
		fileName = filepath.Join(dir, fileName)
		fp, err := os.OpenFile(fileName, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0755)
		if err != nil {
			log.Fatal(err)
		}

		if logger.fp != nil {
			logger.fp.Close()
		}

		logger.fp = fp
		logger.lastEditDay = ymd
		log.SetOutput(fp)
	}

	//t := fmt.Sprintf("%04d%02d%02d %02d:%02d:%02d", year,month,day, now.Hour(), now.Minute(), now.Second())
	//log.SetPrefix(t)
	switch l {
	case LevelTrace:
		log.Printf("[trace]: %s",  fmt.Sprint(v...))
	case LevelDebug:
		log.Printf("[debug]: %s",  fmt.Sprint(v...))
	case LevelInfo:
		log.Printf("[info]: %s", fmt.Sprint(v...))
	case LevelWarn:
		log.Printf("[warn]: %s",  fmt.Sprint(v...))
	case LevelError:
		log.Printf("[error]: %s", fmt.Sprint(v...))
	case LevelFatal:
		log.Fatalf("[fatal]: %s", fmt.Sprint(v...))
	}
}


var std = new(Logger)

func Trace(v ...interface{})  {
	std.Trace(v...)
}

func Debug(v ...interface{})  {
	std.Debug(v...)
}

func Info(v ...interface{})  {
	std.Info(v...)
}

func Warn(v ...interface{})  {
	std.Warn(v...)
}

func Error(v ...interface{})  {
	std.Error(v...)
}

func Fatal(v ...interface{})  {
	std.Fatal(v...)
}

func SetFile(fileName string)  {
	std.SetFile(fileName)
}

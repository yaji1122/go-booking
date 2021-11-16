package logger

import (
	"log"
	"os"
)

type Logger struct {
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
}

func CreateLogger() *Logger {
	logInstance := Logger{
		InfoLogger:  log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
		ErrorLogger: log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
	}
	return &logInstance
}

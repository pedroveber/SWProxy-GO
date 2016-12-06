package log

import (
	"io"
	stdLog "log"
	"os"
)

type Logger struct {
	// write to file only
	debugLogger *stdLog.Logger
	// write to both file and stdout
	infoLogger *stdLog.Logger
}

var logger *Logger

func init() {
	logf, err := os.OpenFile("log.log", os.O_WRONLY|os.O_CREATE, 0640)
	if err != nil {
		stdLog.Fatalln(err)
	}
	multiWriter := io.MultiWriter(logf, os.Stdout)
	dLogger := stdLog.New(logf, "", stdLog.LstdFlags)
	iLogger := stdLog.New(multiWriter, "", stdLog.LstdFlags)
	logger = &Logger{dLogger, iLogger}
}

func (l *Logger) Print(v ...interface{}) {
	l.infoLogger.Println(v...)
}

func (l *Logger) Println(v ...interface{}) {
	l.infoLogger.Println(v...)
}

func (l *Logger) Printf(format string, v ...interface{}) {
	l.infoLogger.Printf(format, v...)
}

func (l *Logger) Debug(v ...interface{}) {
	l.debugLogger.Print(v...)
}

func (l *Logger) Debugln(v ...interface{}) {
	l.debugLogger.Println(v...)
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	l.debugLogger.Printf(format, v...)
}

func (l *Logger) Fatal(v ...interface{}) {
	l.infoLogger.Fatal(v...)
}

func (l *Logger) Fatalln(v ...interface{}) {
	l.infoLogger.Fatalln(v...)
}

func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.infoLogger.Fatalf(format, v...)
}

func Print(v ...interface{}) {
	logger.Print(v...)
}

func Println(v ...interface{}) {
	logger.Println(v...)
}

func Printf(format string, v ...interface{}) {
	logger.Printf(format, v...)
}

func Debug(v ...interface{}) {
	logger.Debug(v...)
}

func Debugln(v ...interface{}) {
	logger.Debugln(v...)
}

func Debugf(format string, v ...interface{}) {
	logger.Debugf(format, v...)
}

func Fatal(v ...interface{}) {
	logger.Fatal(v...)
}

func Fatalln(v ...interface{}) {
	logger.Fatalln(v...)
}

func Fatalf(format string, v ...interface{}) {
	logger.Fatalf(format, v...)
}

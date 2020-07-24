package LogUtil

import (
	"log"
	"os"
	"runtime"
	"sync"
	"time"
)

type Logger struct {
	logger         *log.Logger
	FilePrefixName string
	rwLock         *sync.RWMutex
}

var StdLogger = &Logger{
	logger: log.New(os.Stdout, "", log.LstdFlags),
	rwLock: new(sync.RWMutex),
}

func NewLogger(prefix string) *Logger {
	l := Logger{
		FilePrefixName: prefix,
	}
	l.rwLock = new(sync.RWMutex)
	now := time.Now()
	year, month, day := now.Date()
	today := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	filename := prefix + "_" + today.Format("20060102") + ".log"
	file, err := os.Create(filename)
	if err != nil {
	}
	l.logger = log.New(file, "", log.LstdFlags)
	go func() {
		today = today.Add(24 * time.Hour)
		t := time.NewTimer(today.Sub(time.Now()))
		for range t.C {
			func() {
				filename = prefix + "_" + today.Format("20060102") + ".log"
				l.rwLock.Lock()
				_ = file.Close()
				defer l.rwLock.Unlock()
				file, err := os.Create(filename)
				if err != nil {
					return
				}
				l.logger = log.New(file, "", log.LstdFlags)
			}()
		}
	}()
	return &l
}
func (l *Logger) Info(m string) {
	file, line := getFileAndLine()
	l.rwLock.RLock()
	defer l.rwLock.RUnlock()
	l.logger.Printf("INFO: %s\tfile:%s line:%d\n", m, file, line)
}
func (l *Logger) Error(m string) {
	file, line := getFileAndLine()
	l.rwLock.RLock()
	defer l.rwLock.RUnlock()
	l.logger.Printf("ERROR: %s\tfile:%s line:%d\n", m, file, line)
}
func (l *Logger) Fatal(m string) {
	file, line := getFileAndLine()
	l.rwLock.RLock()
	defer l.rwLock.RUnlock()
	l.logger.Fatalf("FATAL: %s\tfile:%s line:%d\n", m, file, line)
}
func getFileAndLine() (string, int) {
	_, file, line, ok := runtime.Caller(2)
	if ok {
		return file, line
	} else {
		return "", 0
	}
}

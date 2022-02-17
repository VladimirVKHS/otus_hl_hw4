package logger

import (
	"log"
	"os"
	"sync"
)

var infoLoggers []*log.Logger
var errorLoggers []*log.Logger
var f *os.File
var mu sync.RWMutex = sync.RWMutex{}

func Init() {
	var err error
	f, err = os.OpenFile("logs/log.txt", os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		panic("Logger error!")
	}
	infoLoggers = []*log.Logger{
		log.New(f, "INFO\t", log.Ldate|log.Ltime),
		log.New(os.Stdout, "INFO\t", log.Ltime),
	}
	errorLoggers = []*log.Logger{
		log.New(f, "ERROR\t", log.Ldate|log.Ltime),
		log.New(os.Stderr, "ERROR\t", log.Ltime),
	}
}

func Close() {
	if f != nil {
		f.Close()
	}
}

func Info(message string) {
	mu.Lock()
	defer mu.Unlock()
	for _, l := range infoLoggers {
		l.Println("\t" + message)
	}
}

func Error(message string) {
	for _, l := range errorLoggers {
		l.Println("\t" + message)
	}
}

package logger

import (
	"log"
	"os"
)

var (
	Debug *log.Logger
	Error *log.Logger
	Info  *log.Logger
	Warn  *log.Logger
)

func CreateLogger() error {
	makeLogsDir()
	debugFile, err := os.OpenFile("logs/debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		log.Fatal("Cannot access DEBUG log file:", err)
	}
	errorFile, err := os.OpenFile("logs/error.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		log.Fatal("Cannot access ERROR log file:", err)
	}
	infoFile, err := os.OpenFile("logs/info.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		log.Fatal("Cannot access INFO log file:", err)
	}
	warnFile, err := os.OpenFile("logs/warn.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		log.Fatal("Cannot access WARN log file:", err)
	}
	Debug = log.New(debugFile, "DEBUG:", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(errorFile, "ERROR:", log.Ldate|log.Ltime|log.Lshortfile)
	Info = log.New(infoFile, "INFO:", log.Ldate|log.Ltime|log.Lshortfile)
	Warn = log.New(warnFile, "INFO:", log.Ldate|log.Ltime|log.Lshortfile)
	return nil
}

func makeLogsDir() {
	if _, err := os.Stat("./logs"); err != nil {
		if os.IsNotExist(err) {
			os.Mkdir("./logs", 0777)
		}
	}
}

package config

import (
	"io"
	"log"
	"os"
)

var (
	TransactionLogger *log.Logger
	LogFile           *os.File
)

func InitLogger(cfg *EnvConfig) {
	var err error
	LogFile, err = os.OpenFile(cfg.LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	multi := io.MultiWriter(os.Stdout, LogFile)
	TransactionLogger = log.New(multi, "[TX] ", log.Ldate|log.Ltime|log.Lshortfile)
}

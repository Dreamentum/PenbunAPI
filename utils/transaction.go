// utils/tx.go
package utils

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"time"
)

var txLogger *log.Logger

func init() {
	// Determine the transaction log path. Prefer TRANSACTION_LOG_FILE, then LOG_FILE, then a repo default.
	logFile := os.Getenv("TRANSACTION_LOG_FILE")
	if logFile == "" {
		logFile = os.Getenv("LOG_FILE")
	}
	if logFile == "" {
		logFile = "logs/transaction.log"
	}

	// Ensure the directory exists (create it if needed)
	dir := filepath.Dir(logFile)
	if dir != "." && dir != "/" && dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Printf("[WARN] failed creating transaction log directory %s: %v", dir, err)
		}
	}

	f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		// ถ้าเปิดไม่ได้ ให้ fallback ไป stdout แต่ไม่ทำให้แอปดับ
		log.Printf("[WARN] open transaction log failed: %v", err)
		txLogger = log.New(os.Stdout, "", log.LstdFlags)
	} else {
		txLogger = log.New(f, "", log.LstdFlags)
	}
}

func ExecuteTransaction(db *sql.DB, steps []func(tx *sql.Tx) error) error {
	start := time.Now()
	tx, err := db.Begin()
	if err != nil {
		txLogger.Printf("[BEGIN][ERR] %v", err)
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			txLogger.Printf("[PANIC][ROLLBACK] %v", p)
			panic(p)
		}
	}()

	for i, step := range steps {
		if err := step(tx); err != nil {
			_ = tx.Rollback()
			txLogger.Printf("[ROLLBACK] step=%d err=%v elapsed=%s", i+1, err, time.Since(start))
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		txLogger.Printf("[COMMIT][ERR] %v elapsed=%s", err, time.Since(start))
		return err
	}

	txLogger.Printf("[COMMIT][OK] steps=%d elapsed=%s", len(steps), time.Since(start))
	return nil
}

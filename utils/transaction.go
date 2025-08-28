// utils/tx.go
package utils

import (
	"database/sql"
	"log"
	"os"
	"time"
)

var txLogger *log.Logger

func init() {
	// ให้แน่ใจว่าโฟลเดอร์มีอยู่
	_ = os.MkdirAll("/logs/transactions", 0755)

	f, err := os.OpenFile("/logs/transactions/tx.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		// ถ้าเปิดไม่ได้ ให้ fallback ไป stdout แต่ไม่ทำให้แอปดับ
		log.Printf("[WARN] open tx log failed: %v", err)
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

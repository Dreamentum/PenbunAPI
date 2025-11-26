package main

import (
	"fmt"
	"os"

	"PenbunAPI/config"
)

func main() {
	os.Setenv("LOG_FILE", "logs/testlogs/transaction.log")
	// Ensure test directory is removed if exists
	_ = os.RemoveAll("logs/testlogs")
	config.InitLogger()

	// Print out whether the file exists
	if _, err := os.Stat("logs/testlogs/transaction.log"); os.IsNotExist(err) {
		fmt.Println("transaction log file not created")
	} else if err != nil {
		fmt.Println("stat error:", err)
	} else {
		fmt.Println("transaction log file created successfully")
	}
}

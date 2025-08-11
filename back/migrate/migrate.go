package main

import (
	"fmt"

	"github.com/yanatoritakuma/budget/back/db"
	"github.com/yanatoritakuma/budget/back/model"
)

func main() {
	dbConn := db.NewDB()
	defer fmt.Println("Successfully Migrated")
	defer db.CloseDB(dbConn)
	dbConn.AutoMigrate(
		&model.User{},
		&model.Expense{},
	)
}

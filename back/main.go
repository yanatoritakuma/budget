package main

import (
	"github.com/yanatoritakuma/budget/back/controller"
	"github.com/yanatoritakuma/budget/back/db"
	"github.com/yanatoritakuma/budget/back/repository"
	"github.com/yanatoritakuma/budget/back/router"
	"github.com/yanatoritakuma/budget/back/usecase"
)

func main() {
	db := db.NewDB()

	// User関連の依存関係
	userRepository := repository.NewUserRepository(db)
	userUsecase := usecase.NewUserUsecase(userRepository)
	userController := controller.NewUserController(userUsecase)

	// Expense関連の依存関係
	expenseRepository := repository.NewExpenseRepository(db)
	expenseUsecase := usecase.NewExpenseUsecase(expenseRepository)
	expenseController := controller.NewExpenseController(expenseUsecase)

	r := router.NewRouter(userController, expenseController)
	r.Run(":8080")
}

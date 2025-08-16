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

	// Repositories
	userRepository := repository.NewUserRepository(db)
	householdRepository := repository.NewHouseholdRepository(db)
	expenseRepository := repository.NewExpenseRepository(db)

	// Usecases
	userUsecase := usecase.NewUserUsecase(userRepository, householdRepository)
	expenseUsecase := usecase.NewExpenseUsecase(expenseRepository)

	// Controllers
	userController := controller.NewUserController(userUsecase)
	expenseController := controller.NewExpenseController(expenseUsecase)

	r := router.NewRouter(userController, expenseController)
	r.Run(":8080")
}

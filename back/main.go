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
	expenseUsecase := usecase.NewExpenseUsecase(expenseRepository, userRepository)
	householdUsecase := usecase.NewHouseholdUsecase(householdRepository, userRepository)

	// Controllers
	userController := controller.NewUserController(userUsecase)
	expenseController := controller.NewExpenseController(expenseUsecase)
	householdController := controller.NewHouseholdController(householdUsecase)

	r := router.NewRouter(userController, expenseController, householdController)
	r.Run(":8080")
}

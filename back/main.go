package main

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
	"github.com/yanatoritakuma/budget/back/controller"
	"github.com/yanatoritakuma/budget/back/db"
	"github.com/yanatoritakuma/budget/back/repository"
	"github.com/yanatoritakuma/budget/back/router"
	"github.com/yanatoritakuma/budget/back/usecase"
)

var ginLambda *ginadapter.GinLambda

// setupRouter initializes the database, repositories, usecases, controllers and router.
func setupRouter() *gin.Engine {
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

	return router.NewRouter(userController, expenseController, householdController)
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return ginLambda.ProxyWithContext(ctx, req)
}

func main() {
	r := setupRouter()

	if _, ok := os.LookupEnv("LAMBDA_TASK_ROOT"); ok {
		ginLambda = ginadapter.New(r)
		lambda.Start(Handler)
	} else {
		log.Fatal(r.Run(":8080"))
	}
}

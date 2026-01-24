package main

import (
	"context"
	"log"
	"os"
	"strings"

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

var ginLambda *ginadapter.GinLambdaV2

// setupRouter initializes the database, repositories, usecases, controllers and router.
func setupRouter() *gin.Engine {
	dbInstance := db.NewDB()

	// Repositories
	userRepoImpl := repository.NewUserRepositoryImpl(dbInstance)
	householdRepoImpl := repository.NewHouseholdRepositoryImpl(dbInstance)
	expenseRepository := repository.NewExpenseRepositoryImpl(dbInstance)
	uow := repository.NewUnitOfWork(dbInstance)

	// Usecases
	expenseUsecase := usecase.NewExpenseUsecase(expenseRepository, userRepoImpl)
	userUsecase := usecase.NewUserUsecase(userRepoImpl, householdRepoImpl, uow)

	// Controllers
	expenseController := controller.NewExpenseController(expenseUsecase)

	// New router signature
	return router.NewRouter(dbInstance, expenseController, userRepoImpl, householdRepoImpl, uow, userUsecase)
}

func Handler(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	req.RawPath = strings.TrimPrefix(req.RawPath, "/prod")

	// gin-adapter を介してリクエストを処理
	resp, err := ginLambda.ProxyWithContext(ctx, req)
	if err != nil {
		// エラーが発生した場合、ログに出力してエラーレスポンスを返す
		log.Printf("Error proxying request: %v", err)
		return events.APIGatewayV2HTTPResponse{StatusCode: 500}, err
	}

	// API Gateway (HTTP API) は、Set-Cookie ヘッダを `headers` ではなく `cookies` フィールドで受け取る必要がある。
	// gin-adapter が `MultiValueHeaders` に `Set-Cookie` を設定するため、ここから `cookies` フィールドに明示的に移し替える。
	if setCookies, ok := resp.MultiValueHeaders["Set-Cookie"]; ok {
		resp.Cookies = append(resp.Cookies, setCookies...)
		delete(resp.MultiValueHeaders, "Set-Cookie")
	}

	// 念のため、単一のヘッダーも確認
	if setCookie, ok := resp.Headers["Set-Cookie"]; ok {
		resp.Cookies = append(resp.Cookies, setCookie)
		delete(resp.Headers, "Set-Cookie")
	}

	return resp, nil
}

func main() {
	r := setupRouter()

	if _, ok := os.LookupEnv("LAMBDA_TASK_ROOT"); ok {
		ginLambda = ginadapter.NewV2(r)
		lambda.Start(Handler)
	} else {
		log.Fatal(r.Run(":8080"))
	}
}

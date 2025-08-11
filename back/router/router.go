package router

import (
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/yanatoritakuma/budget/back/controller"
)

func NewRouter(
	uc controller.IUserController,
	ec controller.IExpenseController,
) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// プロキシの設定
	r.ForwardedByClientIP = true
	r.TrustedPlatform = "X-Forwarded-For" // PlatformGoogleCloudを直接的なヘッダー指定に変更

	// CORSの設定
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", os.Getenv("FE_URL")},
		AllowMethods:     []string{"GET", "PUT", "POST", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-CSRF-Token"},
		AllowCredentials: true,
	}))

	// CSRFの設定
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", os.Getenv("FE_URL"))
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}
		c.Next()
	})

	// CSRF保護を適用
	r.Use(csrfMiddleware(uc))

	// 認証不要なエンドポイント
	r.POST("/signup", gin.HandlerFunc(uc.SignUp))
	r.POST("/login", gin.HandlerFunc(uc.LogIn))
	r.POST("/logout", gin.HandlerFunc(uc.LogOut))
	r.GET("/csrf", gin.HandlerFunc(uc.CsrfToken))

	// 認証が必要なエンドポイント
	auth := r.Group("/user")
	auth.Use(authMiddleware())
	{
		auth.GET("", gin.HandlerFunc(uc.GetLoggedInUser))
		auth.PUT("", gin.HandlerFunc(uc.UpdateUser))
		auth.DELETE("/:userId", gin.HandlerFunc(uc.DeleteUser))
	}

	// 支出管理のエンドポイント（認証必要）
	expenses := r.Group("/expenses")
	expenses.Use(authMiddleware())
	{
		expenses.POST("", gin.HandlerFunc(ec.CreateExpense))
	}

	return r
}

func csrfMiddleware(uc controller.IUserController) gin.HandlerFunc {
	return func(c *gin.Context) {
		// GETリクエストやOPTIONSリクエストはスキップ
		if c.Request.Method == "GET" || c.Request.Method == "OPTIONS" {
			c.Next()
			return
		}

		// CSRFトークンの検証
		token := c.GetHeader("X-CSRF-Token")
		if token == "" {
			c.JSON(http.StatusForbidden, gin.H{"error": "CSRF token missing"})
			c.Abort()
			return
		}

		// セッションIDの取得（この例ではJWTトークンを使用）
		sessionID, err := c.Cookie("token")
		if err != nil {
			sessionID = "default" // フォールバック値
		}

		// トークンの検証
		if !uc.ValidateCSRFToken(sessionID, token) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Invalid CSRF token"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie("token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		token, err := jwt.Parse(cookie, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(os.Getenv("SECRET")), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		c.Set("user", claims)
		// user_idを設定
		if userID, ok := claims["user_id"].(float64); ok {
			c.Set("user_id", uint(userID))
		}
		c.Next()
	}
}

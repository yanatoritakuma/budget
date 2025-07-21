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

	return r
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
		c.Next()
	}
}

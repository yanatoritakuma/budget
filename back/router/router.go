package router

import (
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/yanatoritakuma/budget/back/controller"
	"github.com/yanatoritakuma/budget/back/repository" // Added
	"github.com/yanatoritakuma/budget/back/usecase"    // Added
	"gorm.io/gorm"                                     // Added
)

func NewRouter(
	db *gorm.DB,
	ec controller.IExpenseController,
	hc controller.IHouseholdController,
) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// プロキシの設定
	r.ForwardedByClientIP = true
	r.TrustedPlatform = "X-Forwarded-For"

	// CORSの設定
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", os.Getenv("FE_URL")},
		AllowMethods:     []string{"GET", "PUT", "POST", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-CSRF-Token"},
		AllowCredentials: true,
	}))

	// CSRFヘッダー設定
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", os.Getenv("FE_URL"))
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}
		c.Next()
	})

	// --- Dependency Injection for User module ---
	userRepo := repository.NewUserRepositoryImpl(db)
	// TODO: Replace nil with a proper IHouseholdRepository implementation
	// when the Household domain is also refactored.
	userUsecase := usecase.NewUserUsecase(userRepo, nil)
	userController := controller.NewUserController(userUsecase)
	// --- End Dependency Injection for User module ---

	// CSRF保護を適用
	r.Use(csrfMiddleware(userController)) // Use instantiated userController

	// -------------------------
	// 認証不要ルート
	// -------------------------
	r.POST("/signup", gin.HandlerFunc(userController.SignUp)) // Use instantiated userController
	r.POST("/login", gin.HandlerFunc(userController.LogIn))   // Use instantiated userController
	r.POST("/logout", gin.HandlerFunc(userController.LogOut)) // Use instantiated userController
	r.GET("/csrf", gin.HandlerFunc(userController.CsrfToken)) // Use instantiated userController

	// -------------------------
	// 認証必須ルート
	// -------------------------
	auth := r.Group("/user")
	auth.Use(authMiddleware())
	{
		auth.GET("", gin.HandlerFunc(userController.GetLoggedInUser))       // Use instantiated userController
		auth.PUT("", gin.HandlerFunc(userController.UpdateUser))            // Use instantiated userController
		auth.DELETE("/:userId", gin.HandlerFunc(userController.DeleteUser)) // Use instantiated userController
	}

	// 支出管理のエンドポイント（認証必要）
	expenses := r.Group("/expenses")
	expenses.Use(authMiddleware())
	{
		expenses.POST("", gin.HandlerFunc(ec.CreateExpense))
		expenses.GET("", gin.HandlerFunc(ec.GetExpense))
		expenses.PUT("/:id", gin.HandlerFunc(ec.UpdateExpense))
		expenses.DELETE("/:id", gin.HandlerFunc(ec.DeleteExpense))
	}

	// 世帯管理のエンドポイント（認証必要）
	household := r.Group("/household")
	household.Use(authMiddleware())
	{
		household.GET("/users", gin.HandlerFunc(userController.GetHouseholdUsers)) // Use instantiated userController
		household.POST("/invite-code", gin.HandlerFunc(hc.GenerateInviteCode))
		household.POST("/join", gin.HandlerFunc(userController.JoinHousehold)) // Use instantiated userController
	}

	return r
}

// ==========================
// CSRF Middleware
// ==========================
func csrfMiddleware(uc controller.IUserController) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "GET" || c.Request.Method == "OPTIONS" {
			c.Next()
			return
		}
		token := c.GetHeader("X-CSRF-Token")
		if token == "" {
			c.JSON(http.StatusForbidden, gin.H{"error": "CSRF token missing"})
			c.Abort()
			return
		}
		sessionID, err := c.Cookie("token")
		if err != nil {
			sessionID = "default"
		}
		if !uc.ValidateCSRFToken(sessionID, token) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Invalid CSRF token"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// ==========================
// Auth Middleware
// ==========================
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
		if userID, ok := claims["user_id"].(float64); ok {
			c.Set("user_id", uint(userID))
		}
		c.Next()
	}
}

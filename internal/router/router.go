package router

import (
	"github.com/Dimpal-Kalita/LetsNormalizeIt-backend/config"
	"github.com/Dimpal-Kalita/LetsNormalizeIt-backend/internal/auth"
	"github.com/Dimpal-Kalita/LetsNormalizeIt-backend/internal/handlers"
	"github.com/Dimpal-Kalita/LetsNormalizeIt-backend/internal/middleware"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewRouter(db *mongo.Client, jwtSecret string) *gin.Engine {
	r := gin.Default()

	GinMode := config.Loadconfig().GIN_MODE
	if GinMode == "release" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	r.Use(middleware.CORSmiddleware())

	authService := &auth.AuthHandler{DB: db.Database("LetsNormaliZeIt"), JWTSecret: jwtSecret}
	authHandler := &handlers.AuthHandler{AuthService: authService}

	r.POST("/register", authHandler.Register)
	r.PATCH("/verify/:token", authHandler.VerifyEmail)
	r.GET("/login", authHandler.Login)
	r.PATCH("/forgot-password", authHandler.ForgotPassword)
	r.PATCH("/reset-password/:token", authHandler.ResetPassword)
	r.GET("/validate-token", authHandler.ValidateToken)

	// user functions
	userHandler := &handlers.UserHandler{DB: db.Database("LetsNormaliZeIt"), JWTSecret: jwtSecret}
	r.GET("/userEmail/:token", userHandler.GetEmailFromToken)

	// blog functions
	blogHandler := &handlers.BlogHandler{DB: db.Database("LetsNormaliZeIt")}
	r.POST("/create-blog", blogHandler.CreateBlog)
	r.PATCH("/update-blog/:id", blogHandler.UpdateBlog)
	r.GET("/blogs", blogHandler.GetBlogs)
	r.GET("/blog/:id", blogHandler.GetBlog)
	r.GET("/user/blog", blogHandler.GetUserBlog)
	r.DELETE("/blog/:id", blogHandler.DeleteBlog)
	return r
}

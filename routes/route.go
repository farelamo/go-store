package routes

import (
	"store/controller"
	"store/initializer"
	"store/middleware"
	"store/repository"

	"github.com/gin-gonic/gin"
)

func Routes() {
	mysql, err := initializer.MysqlInit()
	if err != nil {
		panic("Failed to connect to the database: " + err.Error())
	}

	rdb, err := initializer.RedisInit()
	if err != nil {
		panic("Failed to connect to Redis: " + err.Error())
	}

	if err := initializer.MysqlMigrate(mysql); err != nil {
		panic("Failed to migrate the database: " + err.Error())
	}

	// repositories
	productRepo := repository.NewProductRepository(mysql)
	authRepo := repository.NewAuthRepository(rdb)
	userRepo := repository.NewUserRepository(mysql)

	// controllers
	productController := controller.NewProductController(productRepo)
	authController := controller.NewAuthControlller(authRepo, userRepo)

	r := gin.New()
	r.POST("/login", authController.Login)
	r.POST("/register", authController.Register)
	r.POST("/refresh-token", authController.RefreshToken)
	r.POST("/revoke-token", authController.RevokeToken)

	api := r.Group("/api/v1", middleware.CORSMiddleware(), middleware.XSSProtectionMiddleware(), middleware.JWTAuthMiddleware(authRepo))
	{
		product := api.Group("/product")
		{
			product.GET("/", productController.GetProducts)
			product.GET("/:id", productController.GetProductByID)
			product.POST("/", productController.CreateProduct)
		}
	}

	if err := r.Run("0.0.0.0:8082"); err != nil {
		panic("Failed to start the server: " + err.Error())
	}
}

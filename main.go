package main

import (
	"pro1/controller"
	"pro1/database"
	"pro1/middleware"
	"pro1/repository"
	"pro1/service"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

var (
	db             *gorm.DB                  = database.GetDataBase()
	userRepository repository.UserRepository = repository.NewUserRepository(db)
	dataRepository repository.DataRepository = repository.NewDataRepository(db)
	jwtService     service.JWTService        = service.NewJWTService()
	userService    service.UserService       = service.NewUserService(userRepository)
	dataService    service.DataService       = service.NewDataService(dataRepository)
	authService    service.AuthService       = service.NewAuthService(userRepository)
	authController controller.AuthController = controller.NewAuthController(authService, jwtService)
	userController controller.UserController = controller.NewUserService(userService, jwtService)
	dataController controller.DataController = controller.NewDataController(dataService, jwtService)
)


// Негізгі паток main, егер жаңа запрос келіп түссе оны жаңа басқа патокка салып жібереді(Бұл дегеніміз әп бір запросқа жаңа паток яғни бұғаттамайды)
func main() {
	defer database.CloseDataBase(db)
	r := gin.Default()

	auth := r.Group("/api/auth")
	{
		auth.POST("/register", authController.Register)
		auth.POST("/login", authController.Login)
	}

	userRoutes := r.Group("api/user", middleware.AuthorizeJWT(jwtService))
	{
		userRoutes.GET("/profile", userController.Profile)
		userRoutes.PUT("/profile", userController.Update)
	}

	data := r.Group("/api/data", middleware.AuthorizeJWT(jwtService))
	{
		data.GET("/", dataController.All)
		data.GET("/qu/:q/covid", dataController.FindByString)
		data.GET("/:id", dataController.FindById)
		data.POST("/", dataController.Insert)
		data.PUT("/:id", dataController.Update)
		data.DELETE("/:id", dataController.Delete)
	}

	r.Run(":8080")
}

package route

import (
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// gin-swagger middleware
// swagger embed files


func InitRoute(api *gin.Engine) {
	apiGroup := api.Group("/api")
	v1 := apiGroup.Group("/v1")

	v1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}

// func userRoute(r *gin.RouterGroup) {
// 	userRepository := repository.NewUserRepository()
// 	h := userPre.NewHandler(
// 		usecase.NewCreateUserUsecase(userRepository),
// 		usecase.NewGetUserByIDUsecase(userRepository),
// 	)
// 	group := r.Group("/users")
// 	group.GET("/:id", h.GetUserByID)
// 	group.POST("", h.CreateUser)
// }
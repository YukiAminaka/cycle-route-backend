package route

import (
	"github.com/YukiAminaka/cycle-route-backend/internal/infrastructure/database/dbgen"
	"github.com/YukiAminaka/cycle-route-backend/internal/infrastructure/repository"
	userPre "github.com/YukiAminaka/cycle-route-backend/internal/presentation/user"
	usecase "github.com/YukiAminaka/cycle-route-backend/internal/usecase/user"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)



func InitRoute(api *gin.Engine, q *dbgen.Queries) {
	apiGroup := api.Group("/api")
	v1 := apiGroup.Group("/v1")

	v1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	{
		userRoute(v1, q)
	}
}

func userRoute(r *gin.RouterGroup, q *dbgen.Queries) {
	userRepository := repository.NewUserRepository(q)
	h := userPre.NewHandler(
		usecase.NewCreateUserUsecase(userRepository),
		usecase.NewGetUserByIDUsecase(userRepository),
	)
	group := r.Group("/users")
	group.GET("/:id", h.GetUserByID)
	group.POST("", h.CreateUser)
}
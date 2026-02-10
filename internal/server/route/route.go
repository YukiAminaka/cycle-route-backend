package route

import (
	"github.com/YukiAminaka/cycle-route-backend/internal/infrastructure/database/dbgen"
	"github.com/YukiAminaka/cycle-route-backend/internal/infrastructure/repository"
	"github.com/YukiAminaka/cycle-route-backend/internal/presentation/middleware"
	routePre "github.com/YukiAminaka/cycle-route-backend/internal/presentation/route"
	userPre "github.com/YukiAminaka/cycle-route-backend/internal/presentation/user"
	routeUsecase "github.com/YukiAminaka/cycle-route-backend/internal/usecase/route"
	userUsecase "github.com/YukiAminaka/cycle-route-backend/internal/usecase/user"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitRoute(api *gin.Engine, q *dbgen.Queries, pool *pgxpool.Pool) {
	k := middleware.NewMiddleware()
	
	apiGroup := api.Group("/api")
	v1 := apiGroup.Group("/v1")

	v1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	{
		userRoute(v1, q, k)
		routeRoute(v1, q, pool, k)
	}
}

func userRoute(r *gin.RouterGroup, q *dbgen.Queries, k *middleware.KratosMiddleware) {
	userRepository := repository.NewUserRepository(q)
	h := userPre.NewHandler(
		userUsecase.NewCreateUserUsecase(userRepository),
		userUsecase.NewGetUserByIDUsecase(userRepository),
	)
	group := r.Group("/users")
	group.GET("/:id",k.Session(), h.GetUserByID)
	group.POST("", h.CreateUser)
}

func routeRoute(r *gin.RouterGroup, q *dbgen.Queries, pool *pgxpool.Pool, k *middleware.KratosMiddleware) {
	routeRepository := repository.NewRouteRepository(q)
	userRepository := repository.NewUserRepository(q)
	txManager := repository.NewTransactionManager(q, pool)

	h := routePre.NewHandler(
		routeUsecase.NewCreateRouteUsecase(userRepository, txManager),
		routeUsecase.NewGetRouteUsecase(routeRepository, userRepository),
		routeUsecase.NewUpdateRouteUsecase(userRepository, txManager, routeRepository),
		routeUsecase.NewDeleteRouteUsecase(userRepository, txManager, routeRepository),
	)

	group := r.Group("/routes")
	group.POST("", k.Session(), h.CreateRoute)
	group.GET("", k.Session(), h.GetRoutesByUserID) // 認証ユーザーのルート一覧
	group.GET("/:route_id", h.GetRouteByID)
	group.PUT("/:route_id", k.Session(), h.UpdateRoute)
	group.DELETE("/:route_id", k.Session(), h.DeleteRoute)
}

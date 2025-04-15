package routes

import (
	"net/http"

	"github.com/anhhuy1010/DATN-cms-customer/controllers"

	docs "github.com/anhhuy1010/DATN-cms-customer/docs"
	"github.com/anhhuy1010/DATN-cms-customer/middleware"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func RouteInit(engine *gin.Engine) {
	userCtr := new(controllers.UserController)

	engine.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	engine.Use(middleware.Recovery())
	docs.SwaggerInfo.BasePath = "/v1"

	apiV1 := engine.Group("/v1")
	apiV1.Use(middleware.RequestLog())

	// Không có RoleMiddleware ở đây
	// Các route không cần xác thực
	apiV1.POST("/customers/login", userCtr.Login)
	apiV1.POST("/customers/sign", userCtr.SignUp)

	// Các route cần xác thực nằm trong group này
	protected := apiV1.Group("/")
	protected.Use(controllers.RoleMiddleware())
	{
		protected.GET("/customers", userCtr.List)
		protected.GET("/users/:uuid", userCtr.Detail)
		protected.POST("/users", userCtr.Create)
		protected.PUT("/users/:uuid", userCtr.Update)
		protected.PUT("/users/:uuid/update-status", userCtr.UpdateStatus)
		protected.DELETE("/users/:uuid", userCtr.Delete)
	}

	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}

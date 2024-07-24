package routes

import (
	"net/http"

	"github.com/anhhuy1010/cms-user/controllers"

	docs "github.com/anhhuy1010/cms-user/docs"
	"github.com/anhhuy1010/cms-user/middleware"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func RouteInit(engine *gin.Engine) {
	userCtr := new(controllers.UserController)
	customerCtr := new(controllers.CustomerController)

	engine.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Auth Service API")
	})
	engine.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	engine.Use(middleware.Recovery())
	docs.SwaggerInfo.BasePath = "/v1"
	apiV1 := engine.Group("/v1")

	//apiV1.Use(middleware.ValidateHeader())
	// apiV1.Use(middleware.VerifyAuth())
	apiV1.Use(middleware.RequestLog())
	{
		apiV1.POST("/users", userCtr.Create)
		apiV1.GET("/users", userCtr.List)
		apiV1.GET("/users/:uuid", userCtr.Detail)
		apiV1.PUT("/users/:uuid", userCtr.Update)
		apiV1.PUT("/users/:uuid/update-status", userCtr.UpdateStatus)
		apiV1.DELETE("/users/:uuid", userCtr.Delete)
		apiV1.POST("/users/login", userCtr.Login)

		apiV1.POST("/customer/login", customerCtr.LoginCustomer)
		apiV1.POST("/customer/signup", customerCtr.SignUpCustomer)
		// apiV1.POST("/account/login/admin", customerCtr.LoginAdmin)
		// apiV1.POST("/account/signup/admin", customerCtr.SignUpAdmin)

	}
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

}

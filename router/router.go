package router

import (
	"WBABEProject-04/controller"
	"WBABEProject-04/docs"
	"WBABEProject-04/logger"
	"fmt"

	"github.com/gin-gonic/gin"
	swgFiles "github.com/swaggo/files"
	ginSwg "github.com/swaggo/gin-swagger"
)

type Router struct {
	ct *controller.Controller
}

func NewRouter(ct *controller.Controller) (*Router, error) {
	r := &Router{ct: ct}
	return r, nil
}
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, X-Forwarded-For, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
func liteAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c == nil {
			c.Abort()
			return
		}
		auth := c.GetHeader("Authorization")

		fmt.Println("Authorization-word", auth)
		c.Next()
	}
}

func (r *Router) Index() *gin.Engine {
	gin.SetMode(gin.DebugMode)
	e := gin.Default()
	e.Use(logger.GinLogger())
	e.Use(logger.GinRecovery(true))
	e.Use(CORS())
	e.GET("/swagger/:any", ginSwg.WrapHandler(swgFiles.Handler))
	docs.SwaggerInfo.Host = "localhost" //swagger 정보 등록
	logger.Info("start server")

	menu := e.Group("/menu", liteAuth())
	{

		// 메뉴 등록
		menu.POST("", r.ct.RegisterMenu)
		menu.DELETE("/:menu", r.ct.DelMenu)
		// menu.GET("/:name", r.ct.GetMenuWithName)
		// menu.GET("/", r.ct.GetMenu)
		// menu.PUT("/", r.ct.UpdateMenu)
	}
	// order := e.Group("menu/order", liteAuth())
	// {
	// 	order.GET("/:name", r.ct.GetMenuWithName)

	// }

	return e
}

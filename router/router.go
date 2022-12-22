package router

import (
	"WBABEProject-04/controller"
	"WBABEProject-04/logger"
	"fmt"

	"github.com/gin-gonic/gin"
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
	order := e.Group("order", liteAuth())
	{
		order.POST("/menu")
	}

	return e
}

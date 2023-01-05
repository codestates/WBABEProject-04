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
	docs.SwaggerInfo.Host = "localhost"
	logger.Info("start server")

	menu := e.Group("/menu", liteAuth())
	{
		// 메뉴 리스트 가져오기
		menu.GET("", r.ct.GetMenuList)
		// 메뉴 등록
		menu.POST("", r.ct.RegisterMenu)
		// 메뉴 삭제
		menu.DELETE("/:menu", r.ct.DeleteMenu)
		// 메뉴 수정
		menu.PUT("/:menu", r.ct.UpdateMenu)
		// 한가지 메뉴를 가져온다
		menu.GET("/:menu", r.ct.GetOneMenu)
		// 주문 요청 들어온 내역을 보여준다.
		menu.GET("/order", r.ct.GetOrderList)
		// 주문 요청으로 들어온 내역의 상태를 변경한다.
		menu.PUT("/order", r.ct.UpdateOrderStatus)

	}
	order := e.Group("/order", liteAuth())
	{
		// 주문 요청
		order.POST("", r.ct.OrderMenu)
		// 메뉴를 정렬해서 가져온다.
		order.GET("", r.ct.GetSortedMenu)
		order.GET("/:customerID", r.ct.GetOrderInfo)
		// 주문을 변경할 수 있다
		// order.PUT("/:ordernumber", r.ct.UpdateOrder)
		// 과거 주문 내역을 가져올 수 있다.
		order.GET("/history/:customerID", r.ct.GetOrderHistory)
		// order.POST("", r.ct.OrderMenu)
		// 과거 주문 내역에 있는 메뉴의 평점을 작성할 수 있다.
		order.POST("/history/review/:orderNumber", r.ct.WriteReview)
		order.GET("/history/review/:menuID", r.ct.GetReview)

	}
	return e
}

package router

import (
	"WBABEProject-04/controller"

	"github.com/gin-gonic/gin"
)

type Router struct {
	ct *controller.Controller
}

func NewRouter(ct *controller.Controller) (*Router, error) {
	r := &Router{ct: ct}
	return r, nil
}
func (r *Router) Index() *gin.Engine {
	gin.SetMode(gin.DebugMode)
	e := gin.Default()
	e.GET("/")
	return e
}

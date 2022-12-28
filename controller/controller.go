package controller

import (
	"WBABEProject-04/logger"
	"WBABEProject-04/model"
	"encoding/json"
	"net/http"

	"fmt"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	md *model.Model
}

func NewController(rep *model.Model) (*Controller, error) {
	r := &Controller{md: rep}
	return r, nil
}

// 에러 처리 함수
func (p *Controller) RespError(c *gin.Context, body interface{}, status int, err ...interface{}) {
	logger.Debug("RespError")
	bytes, _ := json.Marshal(body)
	// 사용자에게 전달받은 Path, 전달받은 body, 상태코드, err 메시지
	fmt.Println("Request error", "path", c.FullPath(), "body", bytes, "status", status, "error", err)
	// 클라이언트에게 전달
	c.JSON(status, gin.H{
		// 에러 메시지
		"Error": "Request Error",
		// 경로
		"path": c.FullPath(),
		// body
		"body": bytes,
		// 에러 코드
		"status": status,
		// 에러 객체
		"error": err,
	})
	c.Abort()
}

// GetMenuList godoc
// @Summary call GetMenuList, return ok by json.
// @Description 등록된 메뉴 리스트를 가져온다.
// @Router /menu [get]
func (p *Controller) GetMenuList(c *gin.Context) {
	result, err := p.md.GetMenu()
	if err != nil {
		p.RespError(c, nil, 400, "fail, Not Found Param", nil)
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"res":  "ok",
		"data": result,
	})
}

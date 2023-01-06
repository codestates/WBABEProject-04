package controller

import (
	"WBABEProject-04/logger"
	"WBABEProject-04/model"
	"encoding/json"
	"log"
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

func (p *Controller) RespError(c *gin.Context, body interface{}, status int, err ...interface{}) {
	logger.Debug("RespError")
	bytes, _ := json.Marshal(body)
	fmt.Println("Request error", "path", c.FullPath(), "body", bytes, "status", status, "error", err)
	c.JSON(status, gin.H{
		"Error":  "Request Error",
		"path":   c.FullPath(),
		"body":   bytes,
		"status": status,
		"error":  err,
	})
	c.Abort()
}

// GetMenuList godoc
// @Summary 등록된 메뉴 리스트를 가져옵니다.
// @Description 등록된 메뉴를 JSON 형태로 가져옵니다.
// @Router /menu [get]
func (p *Controller) GetMenuList(c *gin.Context) {
	result, err := p.md.GetMenu("menu")
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

func (p *Controller) CheckMenuInDB(menus []model.Menu) ([]model.Menu, error) {
	resultMenu := []model.Menu{}
	for _, menu := range menus {
		if tempmenu, err := p.md.GetOneMenu("name", menu.Name); err != nil {
			log.Println("fail, Can't find menu.")
			return resultMenu, fmt.Errorf("fail, can't find menu")
		} else {
			if err := p.md.IncreaseMenuVolume(menu); err != nil {
				log.Println("fail, The number of menu cannot be increased.")
				return resultMenu, fmt.Errorf("fail, the number of menu cannot be increased")
			}
			resultMenu = append(resultMenu, tempmenu)
		}
	}
	return resultMenu, nil
}

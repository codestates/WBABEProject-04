package controller

import (
	"WBABEProject-04/logger"
	"WBABEProject-04/model"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	md *model.Model
}

func NewController(rep *model.Model) (*Controller, error) {

	r := &Controller{md: rep}
	return r, nil
}
func (p *Controller) RegisterMenu(c *gin.Context) {
	logger.Debug("RegisterMenu")
	name := c.PostForm("name")         // (필수)
	order := c.PostForm("order")       // 주문 가능 여부 (Default: false)
	quantity := c.PostForm("quantity") // 주문가능 개수(Default: infinity)
	origin := c.PostForm("origin")     // 원산지 (필수, Default: 국내산)
	price := c.PostForm("price")       // 가격 (필수)
	spicy := c.PostForm("spicy")       // 맵기 (Default: normal)

	// 가격을 입력하지 않으면 0원 -> 무료로 생각
	nPrice, err := strconv.Atoi(price)
	if err != nil {
		nPrice = 0
	}

	nQuantity, err := strconv.Atoi(quantity)
	if err != nil {
		// 수량을 따로 정해두지 않음
		nQuantity = 1000000
	}

	// 필수 정보로 들어가야할 정보가 없으면
	if len(name) <= 0 || len(origin) <= 0 {
		p.RespError(c, nil, http.StatusUnprocessableEntity, "parameter not found")
	}

	if len(origin) <= 0 {
		origin = "국내산"
	}

	// 주문 가능 여부 : (Default: false)
	bOrder, err := strconv.ParseBool(order)
	if err != nil {
		bOrder = false
	}

	menu, _ := p.md.GetOneMenu("name", name)
	fmt.Println("menu")
	// 이미 등록된 메뉴가 있으면... 에러
	if menu != (model.Menu{}) {
		p.RespError(c, nil, http.StatusUnprocessableEntity, "already resistery menu", nil)
		return
	}

	// 입력된 정보로 메뉴 생성
	req := model.Menu{Name: name, Order: bOrder, Quantity: nQuantity, Spicy: spicy, Origin: origin, Price: nPrice}

	if err := p.md.CreateMenu(req); err != nil {
		p.RespError(c, nil, http.StatusUnprocessableEntity, "parameter not found", err)
		return
	}
	// 요청 성공 시
	c.JSON(http.StatusOK, gin.H{
		"result": "ok",
	})
	c.Next()
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

// DelMenu godoc
// @Summary call DelMenu, return ok by json.
// @Description 메뉴의 이름을 파라미터로 받아 해당 메뉴를 삭제하는 기능
// @Router /order/menu/:name [delete]
func (p *Controller) DelMenu(c *gin.Context) {
	logger.Debug("DelMenu")
	smenu := c.Param("menu")
	if len(smenu) <= 0 {
		p.RespError(c, nil, http.StatusUnprocessableEntity, "parameter not found", nil)
		return
	}

	_, err := p.md.GetOneMenu("mune", smenu)
	if err != nil {
		p.RespError(c, nil, http.StatusUnprocessableEntity, "exist resistery person", nil)
		return
	}

	if err := p.md.DeleteMenu(smenu); err != nil {
		p.RespError(c, nil, http.StatusUnprocessableEntity, "fail delete db", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": "ok",
	})
	c.Next()
}

// GetMenuWithName godoc
// @Summary call GetMenuWithName, return ok by json.
// @Description 메뉴의 이름을 파라미터로 받아 해당 메뉴의 정보를 가져오는 기능
// @Router /order/menu/:name [get]
func (p *Controller) GetMenuWithName(c *gin.Context) {
	logger.Debug("GetMenuWithName")
	sName := c.Param("name")
	if len(sName) <= 0 {
		p.RespError(c, nil, 400, "fail, Not Found Param", nil)
		c.Abort()
		return
	}
	if per, err := p.md.GetOneMenu("name", sName); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"res":  "fail",
			"body": err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"res":  "ok",
			"body": per,
		})
	}
}

// GetMenu godoc
// @Summary call GetMenu, return ok by json.
// @Description 등록된 메뉴 전체의 리스트를 가져올 수 있다.
// @Router /order/menu [get]
func (p *Controller) GetMenu(c *gin.Context) {
	result := p.md.GetMenuList()
	c.JSON(http.StatusOK, gin.H{
		"res":  "ok",
		"data": result,
	})
}

// GetMenu godoc
// @Summary call GetMenu, return ok by json.
// @Description 등록된 메뉴 전체의 리스트를 가져올 수 있다.
// @Router /order/menu [get]
func (p *Controller) UpdateMenu(c *gin.Context) {

}

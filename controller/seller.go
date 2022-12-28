package controller

import (
	"WBABEProject-04/logger"
	"WBABEProject-04/model"
	"fmt"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// 판매자

// RegisterMenu godoc
// @Summary call RegisterMenu, return ok by json.
// @Description 메뉴의 정보를 JSON으로 입력받아 등록한다.
// @Accept  json
// @Produce  json
// @Router /menu [post]
func (p *Controller) RegisterMenu(c *gin.Context) {
	pMenu := model.Menu{}

	if err := c.ShouldBindJSON(&pMenu); err != nil {
		p.RespError(c, nil, 400, "fail, Not Found Param", nil)
		c.Abort()
		return
	}

	// Name이나 Origin과 같은 필드는 == 0으로 처리할 수 있을 것 같아 보입니다.
	if len(pMenu.Name) <= 0 {
		p.RespError(c, nil, 400, "fail, input your name, please", nil)
		c.Abort()
		return
	}
	if len(pMenu.Origin) <= 0 {
		pMenu.Origin = "국내산"
	}

	pMenu.Review = nil
	menu, _ := p.md.GetOneMenuWithName(pMenu.Name)

	result := reflect.DeepEqual(menu, model.Menu{})
	// 이미 등록된 메뉴가 있으면
	if result != true {
		p.RespError(c, nil, http.StatusUnprocessableEntity, "already resistery menu", nil)
		return
	}

	// 아래의 코드를 short syntax로 처리할 수 있을까요?
	err := p.md.CreateMenu(pMenu)
	if err != nil {
		p.RespError(c, nil, http.StatusUnprocessableEntity, "cannot create.", err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"result": "ok",
	})
	c.Next()
}

// DelMenu godoc
// @Summary call DelMenu, return ok by json.
// @Description 메뉴의 아이디를 파라미터로 받아 해당 메뉴를 삭제하는 기능
// @Router /menu/:id [delete]
func (p *Controller) DelMenu(c *gin.Context) {
	logger.Debug("DelMenu")
	id := c.Param("id")

	pId, err := primitive.ObjectIDFromHex(id)

	_, err = p.md.GetOneMenu(pId)
	if err != nil {
		p.RespError(c, nil, http.StatusUnprocessableEntity, "It is not a registered menu", nil)
		return
	}

	if err := p.md.DeleteMenu(pId); err != nil {
		p.RespError(c, nil, http.StatusUnprocessableEntity, "fail delete db", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": "ok",
	})
	c.Next()
}

// UpdateMenu godoc
// @Summary call UpdateMenu, return ok by json.
// @Description 메뉴의 아이디를 파라미터로 받고 JSON으로 수정하려는 내용을 받아 기존 메뉴의 정보를 변경할 수 있다.
// @Router /menu/:id [put]
func (p *Controller) UpdateMenu(c *gin.Context) {
	id := c.Param("id")
	var recvMenu model.Menu
	err := c.ShouldBindJSON(&recvMenu)
	if err != nil {
		p.RespError(c, nil, 400, "fail, Please enter your json correctly", nil)
		c.Abort()
		return
	}

	pId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		p.RespError(c, nil, http.StatusBadRequest, "ID value could not be verified.", err)
		return
	}
	recvMenu.MenuId = pId
	menu, err := p.md.GetOneMenu(recvMenu.MenuId)
	if err != nil {
		p.RespError(c, nil, http.StatusUnprocessableEntity, "fail update db", err)
		return
	}
	var tempMenu model.Menu
	tempMenu.Name = menu.Name

	// 주문 가능 여부가 true면 true
	if recvMenu.Available {
		tempMenu.Available = true
	} else {
		tempMenu.Available = menu.Available
	}
	// 입력 받은 수량이 0이 아니면 (입력을 받은 것이므로)
	if recvMenu.Quantity != 0 {
		tempMenu.Quantity = recvMenu.Quantity
	}

	// 원산지를 입력 받으면 수정
	if len(recvMenu.Origin) > 0 {
		tempMenu.Origin = recvMenu.Origin
	} else {
		tempMenu.Origin = menu.Origin
	}

	// Price를 입력 받으면 수정
	if recvMenu.Price > 0 {
		tempMenu.Price = recvMenu.Price
	} else {
		tempMenu.Price = menu.Price
	}

	// Spiciness를 입력 받으면 수정
	if recvMenu.Spiciness > 0 {
		tempMenu.Spiciness = recvMenu.Spiciness
	} else {
		tempMenu.Spiciness = menu.Spiciness
	}
	// 추천 메뉴가 선택되면
	if recvMenu.Favorites {
		tempMenu.Favorites = true
	}

	if err := p.md.UpdateMenu(tempMenu); err != nil {
		p.RespError(c, nil, http.StatusUnprocessableEntity, "parameter not found", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"result": "ok",
	})
}

// 주문 상태 조회(접수된 것 중 접수, 조리중, 배달중, 배달완료 조회가능)
func (p *Controller) GetOrderList(c *gin.Context) {
	var orderList []model.Order
	orderList, err := p.md.GetOrders()
	if err != nil {
		p.RespError(c, nil, http.StatusUnprocessableEntity, "parameter not found", err)
		return
	}
	fmt.Println(orderList)
}

// 주문 상태 변경(접수된 것을 조리중, 배달중, 배달완료)
func (p *Controller) UpdateOrderStatus(c *gin.Context) {

	// 상태 변경하려는 주문 가져오기(접수된 것 이상으로)
	orderList, err := p.md.GetOrders()
	if err != nil {
		p.RespError(c, nil, http.StatusUnprocessableEntity, "parameter not found", err)
		return
	}
	fmt.Println(orderList)

}

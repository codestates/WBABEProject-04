package controller

import (
	"WBABEProject-04/logger"
	"WBABEProject-04/model"
	"encoding/json"

	"fmt"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	md *model.Model
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

// func (p *Controller) GetMenuWithStore(c *gin.Context) {
// 	logger.Debug("GetMenuWithStore")
// 	var temp model.Store
// 	for _, i := range temp.Menus {
// 		fmt.Println(i)

// 	}

// }

// // GetMenuWithName godoc
// // @Summary call GetMenuWithName, return ok by json.
// // @Description 메뉴의 이름을 파라미터로 받아 해당 메뉴의 정보를 가져오는 기능
// // @Router /order/menu/:name [get]
// func (p *Controller) GetMenuWithName(c *gin.Context) {
// 	logger.Debug("GetMenuWithName")
// 	sName := c.Param("name")
// 	if len(sName) <= 0 {
// 		p.RespError(c, nil, 400, "fail, Not Found Param", nil)
// 		c.Abort()
// 		return
// 	}
// 	if per, err := p.md.GetOneMenu("name", sName); err != nil {
// 		c.JSON(http.StatusOK, gin.H{
// 			"res":  "fail",
// 			"body": err.Error(),
// 		})
// 	} else {
// 		c.JSON(http.StatusOK, gin.H{
// 			"res":  "ok",
// 			"body": per,
// 		})
// 	}
// }

// // GetMenu godoc
// // @Summary call GetMenu, return ok by json.
// // @Description 등록된 메뉴 전체의 리스트를 가져올 수 있다.
// // @Router /order/menu [get]
// func (p *Controller) GetMenu(c *gin.Context) {
// 	result := p.md.GetMenuList()
// 	c.JSON(http.StatusOK, gin.H{
// 		"res":  "ok",
// 		"data": result,
// 	})
// }

// // GetMenu godoc
// // @Summary call GetMenu, return ok by json.
// // @Description 등록된 메뉴 전체의 리스트를 가져올 수 있다.
// // @Router /order/menu [get]
// func (p *Controller) UpdateMenu(c *gin.Context) {
// 	var recvMenu model.Menu
// 	err := c.ShouldBindJSON(&recvMenu)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, err.Error())
// 		return
// 	}
// 	var menu model.Menu
// 	if menu, err = p.md.GetOneMenu("name", recvMenu.Name); err != nil {
// 		p.RespError(c, nil, http.StatusUnprocessableEntity, "fail update db", err)
// 		return
// 	}

// 	var tempMenu model.Menu
// 	tempMenu.Name = menu.Name
// 	// 주문 불가능 여부 (true : 주문 불가능)
// 	if recvMenu.Order {
// 		tempMenu.Order = true
// 	} else {
// 		tempMenu.Order = menu.Order
// 	}

// 	// 원산지가 비어있지 않으면
// 	if recvMenu.Origin != "" {
// 		// Client에게 받은 원산지 저장
// 		tempMenu.Origin = recvMenu.Origin
// 	} else {
// 		// 그대로 저장
// 		tempMenu.Origin = menu.Origin
// 	}

// 	// 메뉴 수량이 -1이면 주문 수량이 없음
// 	if recvMenu.Quantity == -1 {
// 		tempMenu.Quantity = -1
// 	} else if recvMenu.Quantity == 0 {
// 		// 메뉴 수량이 0이면 client에게 값을 받지 않은 것이기 때문에 기본 값을 넣는다.
// 		tempMenu.Quantity = menu.Quantity
// 	} else if recvMenu.Quantity < -1 {
// 		// 주문 수량에 값이 있으면 해당 값을 넣는다.
// 		tempMenu.Quantity = menu.Quantity
// 	} else {
// 		tempMenu.Quantity = recvMenu.Quantity
// 	}

// 	// 가격을 받지 않으면 기존 값을 넣는다.
// 	if recvMenu.Price == 0 {
// 		tempMenu.Price = menu.Price
// 	} else {
// 		// 가격을 받으면 해당 값을 넣는다.
// 		tempMenu.Price = recvMenu.Price
// 	}

// 	if recvMenu.Spicy == "" {
// 		tempMenu.Spicy = menu.Spicy
// 	} else if recvMenu.Spicy == "Spicy" {
// 		tempMenu.Spicy = "Spicy"
// 	} else if recvMenu.Spicy == "Very hot" {
// 		tempMenu.Spicy = "Very hot"
// 	} else {
// 		tempMenu.Spicy = "Normal"
// 	}

// 	if err := p.md.UpdateMenu(tempMenu); err != nil {
// 		p.RespError(c, nil, http.StatusUnprocessableEntity, "parameter not found", err)
// 		return
// 	}
// 	c.JSON(http.StatusOK, gin.H{
// 		"result": "ok",
// 	})
// 	c.Next()
// }

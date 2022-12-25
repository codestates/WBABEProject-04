package controller

import (
	"WBABEProject-04/model"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type anableToOrder struct {
	menuName  string
	available bool
	status    int
}

// 주문자

// 주문이 가능한지 확인한다.
func (p *Controller) CheckMenu(c *gin.Context) {
	customer := model.Customer{}
	err := c.ShouldBindJSON(&customer)
	if err != nil {
		p.RespError(c, nil, 400, "fail, Not Found Param", nil)
		c.Abort()
		return
	}
	anable := []anableToOrder{}
	// 주문 하려는 메뉴가 있는지 확인
	orderList := customer.Orders
	// var menus []model.Menu
	for _, order := range orderList {
		// 해당 주문을 접수중으로 업데이트
		err := p.md.UpdateOrderStatus(order, model.Accepting)
		if err != nil {
			p.RespError(c, nil, http.StatusUnprocessableEntity, "fail update db", err)
			return
		}
		for _, menu := range order.Menus {
			menu, err := p.md.GetOneMenu(menu.MenuId)
			// 메뉴가 없으면 리턴
			if err != nil {
				p.RespError(c, nil, http.StatusUnprocessableEntity, "fail find db", err)
				return
			} else {
				// menus = append(menus, menu)
				temp := anableToOrder{}

				temp.menuName = menu.Name
				if menu.Quantity < 1 || menu.Available {
					temp.available = true
					temp.status = model.Cancellation
					err := p.md.UpdateOrderStatus(order, model.Cancellation)
					if err != nil {
						p.RespError(c, nil, http.StatusUnprocessableEntity, "fail update db", err)
						return
					}
				} else {
					anable = append(anable, temp)
				}
			}
		}
	}
	fmt.Println("anable => ", anable)
	// 주문 수량이 여유 있는지 확인
	if len(anable) <= 0 {
		c.JSON(http.StatusOK, gin.H{
			"result": "메뉴가 매진되었습니다.",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"result": "ok",
		})
	}
}

// 주문 요청
func (p *Controller) OrderMenu(c *gin.Context) {
	customer := model.Customer{}
	err := c.ShouldBindJSON(&customer)
	if err != nil {
		p.RespError(c, nil, 400, "fail, Not Found Param", nil)
		c.Abort()
		return
	}
	tempOrder := []model.Menu{}

	orderList := customer.Orders
	for _, order := range orderList {
		for _, menu := range order.Menus {
			temp, err := p.md.GetOneMenu(menu.MenuId)
			if err != nil {
				p.RespError(c, nil, 400, "fail find db", nil)
				c.Abort()
				return
			}
			fmt.Println("=->", temp)
			// order.Menus = copy(tempOrder, temp)
			tempOrder = append(tempOrder, temp)
		}
	}
	fmt.Println("=>", tempOrder)
}

package controller

import (
	"WBABEProject-04/logger"
	"WBABEProject-04/model"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Order struct {
	OrderID primitive.ObjectID `bson:"orderid" json:"orderid"`
	Status  model.Status       `bson:"status" json:"status"`
}

// RegisterMenu godoc
// @Summary 메뉴를 등록합니다.
// @Description JSON형태로 데이터를 전달받아 메뉴를 생성합니다.
// @Accept  json
// @Router /menu [post]
// @Success 200 {string} string "ok"
// @Failure 400 {string} string "fail, Not Found Param"
// @Failure 400 {string} string "fail, enter a menu name, please"
// @Failure 422 {string} string "already resistery menu"
func (p *Controller) RegisterMenu(c *gin.Context) {
	pMenu := model.Menu{}
	if err := c.ShouldBindJSON(&pMenu); err != nil {
		p.RespError(c, nil, 400, "fail, Not Found Param", nil)
		return
	}
	if len(pMenu.Name) <= 0 {
		p.RespError(c, nil, 400, "fail, enter a menu name, please", nil)
		return
	}
	if len(pMenu.Origin) <= 0 {
		pMenu.Origin = "국내산"
	}
	pMenu.Review = nil
	menu, _ := p.md.GetOneMenu("name", pMenu.Name)
	result := reflect.DeepEqual(menu, model.Menu{})
	if !result {
		p.RespError(c, nil, http.StatusUnprocessableEntity, "already resistery menu", nil)
		return
	}

	err := p.md.CreateMenu(pMenu)
	if err != nil {
		p.RespError(c, nil, http.StatusUnprocessableEntity, "cannot create menu.", err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"result": "ok",
	})
	c.Next()
}

// DeleteMenu godoc
// @Summary 메뉴를 삭제합니다.
// @Description arameter 형태로 메뉴 이름을 받아 해당 메뉴를 삭제합니다.
// @Param name path string true "메뉴를 삭제하기 위함"
// @Accept  json
// @Success 200 {string} string "ok"
// @Failure 400 {string} string "It is not a registered menu"
// @Failure 422 {string} string "fail delete db"
// @Router /menu/:name [delete]
func (p *Controller) DeleteMenu(c *gin.Context) {
	logger.Debug("DeleteMenu")
	menuName := c.Param("menu")
	_, err := p.md.GetOneMenu("name", menuName)
	if err != nil {
		p.RespError(c, nil, http.StatusBadRequest, "It is not a registered menu", nil)
		return
	}

	if err := p.md.DeleteMenu(menuName); err != nil {
		p.RespError(c, nil, http.StatusUnprocessableEntity, "fail delete db", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": "ok",
	})
	c.Next()
}

// UpdateMenu godoc
// @Summary 메뉴의 정보를 수정합니다.
// @Description 메뉴의 이름을 파라미터로 받고 JSON으로 수정하려는 내용을 받아 기존 메뉴의 정보를 변경할 수 있다.
// @Param name path string true "메뉴를 삭제하기 위함"
// @Accept  json
// @Success 200 {string} string "ok"
// @Failure 400 {string} string "fail, Please enter your json correctly"
// @Failure 422 {string} string "Menu is not registered"
// @Failure 422 {string} string "The menu cannot be edited."
// @Router /menu/:menu [put]
func (p *Controller) UpdateMenu(c *gin.Context) {
	menuName := c.Param("menu")

	recvMenu := model.Menu{}
	if err := c.ShouldBindJSON(&recvMenu); err != nil {
		p.RespError(c, nil, http.StatusBadRequest, "fail, Please enter your json correctly", nil)
		return
	}

	if tempMenu, err := p.md.GetOneMenu("name", menuName); err != nil {
		p.RespError(c, nil, http.StatusUnprocessableEntity, "Menu is not registered.", err)
		return
	} else {
		if recvMenu.Price == 0 {
			recvMenu.Price = tempMenu.Price
		}
		if recvMenu.Origin == "" {
			recvMenu.Origin = tempMenu.Origin
		}
		if !recvMenu.Favorites {
			recvMenu.Favorites = tempMenu.Favorites
		}
		if !recvMenu.SoldOut {
			recvMenu.SoldOut = tempMenu.SoldOut
		}
		if err := p.md.UpdateMenu(recvMenu, menuName); err != nil {
			p.RespError(c, nil, http.StatusUnprocessableEntity, "The menu cannot be edited.", err)
			return
		} else {
			c.JSON(http.StatusOK, gin.H{
				"result": "ok",
			})
		}
	}
}

// GetOrderList godoc
// @Summary 접수 완료된 주문들의 리스트를 확인할 수 있다.
// @Description 주문 상태가 "접수완료"인 메뉴들을 확인할 수 있다.
// @Accept  json
// @Success 200 {string} string "ok"
// @Failure 422 {string} string "parameter not found"
// @Router /menu/order [get]
func (p *Controller) GetOrderList(c *gin.Context) {
	if orderList, err := p.md.GetOrders(); err != nil {
		p.RespError(c, nil, http.StatusUnprocessableEntity, "parameter not found", err)
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"result": "ok",
			"data":   orderList,
		})
	}

}

// UpdateOrderStatus godoc
// @Summary 주문들의 상태를 변경할 수 있다.
// @Description 주문 상태가 "접수완료"인 메뉴들의 상태를 변경할 수 있다.
// @Accept  json
// @Success 200 {string} string "ok"
// @Failure 400 {string} string "fail, Please enter your json correctly"
// @Failure 400 {string} string "You entered an incorrect status code."
// @Failure 422 {string} string "fail find order"
// @Failure 422 {string} string "fail update status"
// @Router /menu/order [put]
func (p *Controller) UpdateOrderStatus(c *gin.Context) {
	tempOrder := Order{}
	if err := c.ShouldBindJSON(&tempOrder); err != nil {
		p.RespError(c, nil, http.StatusBadRequest, "fail, Please enter your json correctly", nil)
		c.Abort()
		return
	}
	if order, err := p.md.GetOrderByID(tempOrder.OrderID); err != nil {
		p.RespError(c, nil, http.StatusUnprocessableEntity, "fail find order", err)
		return
	} else {
		if tempOrder.Status <= 7 && tempOrder.Status >= 0 {
			if err := p.md.UpdateOrderStatus(order, int(tempOrder.Status)); err != nil {
				p.RespError(c, nil, http.StatusUnprocessableEntity, "fail update status", err)
				return
			} else {
				c.JSON(http.StatusOK, gin.H{
					"result": "ok",
				})
			}
		} else {
			p.RespError(c, nil, http.StatusBadRequest, "You entered an incorrect status code.", nil)
			return
		}
	}
}

// GetOneMenu godoc
// @Summary 한가지 메뉴에 대한 정보를 얻을 수 있다.
// @Description 메뉴의 이름을 파라미터로 받고 JSON 형태로 해당 메뉴에 대한 정보를 얻을 수 있다.
// @Accept  json
// @Success 200 {string} string "ok"
// @Failure 422 {string} string "It is not a registered menu"
// @Router /menu/:name [get]
func (p *Controller) GetOneMenu(c *gin.Context) {
	menuName := c.Param("menu")
	menu, err := p.md.GetOneMenu("name", menuName)
	if err != nil {
		p.RespError(c, nil, http.StatusUnprocessableEntity, "It is not a registered menu", nil)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"result": "ok",
		"data":   menu,
	})
}

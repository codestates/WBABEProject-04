package controller

import (
	"WBABEProject-04/model"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AddOrder godoc
// @Summary 메뉴를 추가할 수 있습니다.
// @Description Parameter로 주문번호를 입력받고, JSON 형태로 추가하고자 하는 메뉴의 정보를 전달하면 메뉴를 추가할 수 있습니다. 배달 중일 경우 신규 주문으로 요청됩니다.
// @Param orderid path string true "메뉴를 추가하기 위함"
// @Success 200 {string} string "ok"
// @Failure 400 {string} string "fail, Not Found Param"
// @Failure 400 {string} string "fail, Not Found JSON"
// @Failure 422 {string} string "fail, Not Found Order"
// @Failure 422 {string} string "fail, Not Found menu"
// @Failure 422 {string} string "failm increase menu volume"
// @Router /order/:orderid [post]
func (p *Controller) AddOrder(c *gin.Context) {
	sOrderNumber := c.Param("ordernumber")
	var recvParam primitive.ObjectID
	var err error
	if recvParam, err = primitive.ObjectIDFromHex(sOrderNumber); err != nil {
		p.RespError(c, nil, http.StatusBadRequest, "fail, Not Found Param", nil)
		return
	}
	var recvOrder model.Order
	err = c.ShouldBindJSON(&recvOrder)
	if err != nil {
		p.RespError(c, nil, http.StatusBadRequest, "fail, Not Found JSON", nil)
		c.Abort()
		return
	}
	tempOrder := model.Order{}
	if tempOrder, err = p.md.GetOrderByID(recvParam); err != nil {
		p.RespError(c, nil, http.StatusUnprocessableEntity, "fail, Not Found Order", err)
		return
	}
	if tempOrder.Status == model.Cooking || tempOrder.Status == model.InDelivery || tempOrder.Status == model.CompleteDelivery {
		newOrder := model.Order{}
		newOrder.CustomerID = recvOrder.CustomerID
		menus := []model.Menu{}
		for _, menu := range recvOrder.Menus {
			if tempmenu, err := p.md.GetOneMenu("name", menu.Name); err != nil {
				p.RespError(c, nil, http.StatusUnprocessableEntity, "fail, Not Found menu", err)
				return
			} else {
				if err := p.md.IncreaseMenuVolume(menu); err != nil {
					p.RespError(c, nil, http.StatusUnprocessableEntity, "failm increase menu volume", err)
					return
				}
				menus = append(menus, tempmenu)
			}
		}
		newOrder.Menus = menus
		newOrder.Status = model.Accepting
		newOrder.ID = primitive.NewObjectID()
		if err := p.md.CreateOrder(newOrder); err != nil {
			p.RespError(c, nil, http.StatusUnprocessableEntity, "fail, Not Found Param", err)
			return
		} else {
			c.JSON(http.StatusCreated, gin.H{
				"res":      "ok",
				"contents": "Adding a menu failed, requesting a new order.",
				"status":   "new order",
				"Data":     newOrder.ID,
			})
		}
	} else {
		menus := []model.Menu{}
		for _, menu := range recvOrder.Menus {
			if tempmenu, err := p.md.GetOneMenu("name", menu.Name); err != nil {
				p.RespError(c, nil, http.StatusUnprocessableEntity, "fail find menu", err)
				return
			} else {
				if err := p.md.IncreaseMenuVolume(menu); err != nil {
					p.RespError(c, nil, http.StatusUnprocessableEntity, "fail increase menu volume", err)
					return
				}
				menus = append(menus, tempmenu)
			}
		}
		recvOrder.Menus = menus
		tempOrder.Menus = append(tempOrder.Menus, recvOrder.Menus...)
		if err = p.md.UpdateOrder(tempOrder, recvParam); err != nil {
			p.RespError(c, nil, http.StatusUnprocessableEntity, "fail update order", err)
			return
		} else {
			c.JSON(http.StatusCreated, gin.H{
				"res":      "ok",
				"contents": "The menu has been added normally.",
				"status":   "add order",
				"Data":     recvParam,
			})
		}
	}
}

// OrderMenu godoc
// @Summary 주문자가 주문을 요청할 수 있습니다.
// @Description json 형태로 주문을 요청합니다. 요청이 완료되면 주문 번호를 json 형태로 리턴합니다.
// @Success 200 {string} string "ok"
// @Failure 400 {string} string "fail, Not Found Param"
// @Failure 422 {string} string "fail, Not Found Menu"
// @Failure 422 {string} string "fail create order"
// @Failure 422 {string} string "fail create customer"
// @Router /order [post]
func (p *Controller) OrderMenu(c *gin.Context) {
	customer := model.Customer{}
	if err := c.ShouldBindJSON(&customer); err != nil {
		p.RespError(c, nil, http.StatusBadRequest, "fail, Not Found Param", nil)
		c.Abort()
		return
	}
	resultCustomer := model.Customer{}
	resultCustomer.Address = customer.Address
	resultCustomer.Nicname = customer.Nicname
	resultCustomer.Phone = customer.Phone
	resultCustomer.ID = primitive.NewObjectID()
	resultOrders := []model.Order{}
	tempOrder := model.Order{}
	for _, order := range customer.Orders {
		tempOrder.Status = model.Accepting
		if menus, err := p.CheckMenuInDB(order.Menus); err != nil {
			p.RespError(c, nil, http.StatusUnprocessableEntity, "fail, Not Found Menu", err)
			return
		} else {
			tempOrder.Menus = menus
			tempOrder.Status = model.Receipt
			tempOrder.CreatedAt = time.Now()
			tempOrder.ID = primitive.NewObjectID()
			tempOrder.CustomerID = resultCustomer.ID
		}
	}
	if err := p.md.CreateOrder(tempOrder); err != nil {
		p.RespError(c, nil, http.StatusUnprocessableEntity, "fail create order", err)
		return
	}
	resultOrders = append(resultOrders, tempOrder)
	resultCustomer.Orders = resultOrders
	if err := p.md.CreateCustomer(resultCustomer); err != nil {
		p.RespError(c, nil, http.StatusUnprocessableEntity, "fail create customer", err)
		return
	} else {
		tempOrder := model.OrderNumber{}
		for _, order := range resultCustomer.Orders {
			tempOrder.OrderList = append(tempOrder.OrderList, (primitive.ObjectID)(order.ID))
		}
		c.JSON(http.StatusOK, gin.H{
			"res":         "ok",
			"orderNumber": tempOrder.OrderList,
		})
	}
}

// UpdateOrder godoc
// @Summary 메뉴를 변경할 수 있다.
// @Description paramert로 주문ID와 JSON으로 변경하려는 내용을 입력하여 메뉴를 변경할 수 있다. 변경 시 상태가 조리중, 배달중 일경우 실패 알림을 보낸다.
// @Param orderid path string true "메뉴를 변경하기 위함"
// @Success 200 {string} string "ok"
// @Failure 400 {string} string "fail, Not Found Param"
// @Failure 400 {string} string "fail, Not Found JSON"
// @Failure 422 {string} string "fail, Not Found Order"
// @Failure 422 {string} string "fail, Not Found Menu"
// @Success 200 {string} string "The Order cannot be changed."
// @Router /order/:ordernumber [put]
func (p *Controller) UpdateOrder(c *gin.Context) {
	sOrderNumber := c.Param("ordernumber")
	var recvOrdernumber primitive.ObjectID
	var err error
	if recvOrdernumber, err = primitive.ObjectIDFromHex(sOrderNumber); err != nil {
		p.RespError(c, nil, http.StatusBadRequest, "fail, Not Found Param", nil)
		return
	}
	var recvOrder model.Order
	if err = c.ShouldBindJSON(&recvOrder); err != nil {
		p.RespError(c, nil, http.StatusBadRequest, "fail, Not Found JSON", nil)
		c.Abort()
		return
	}
	tempOrder := model.Order{}
	if tempOrder, err = p.md.GetOrderByID(recvOrdernumber); err != nil {
		p.RespError(c, nil, http.StatusUnprocessableEntity, "fail, Not Found Order", err)
		return
	}
	if tempOrder.Status >= model.Cooking {
		p.RespError(c, nil, http.StatusOK, "The Order cannot be changed.", nil)
		return
	}
	resultOrder := model.Order{}
	resultOrder.CreatedAt = time.Now()
	resultOrder.ID = primitive.NewObjectID()
	resultOrder.CreatedAt = tempOrder.CreatedAt
	fmt.Println(recvOrder.Menus)
	if menus, err := p.CheckMenuInDB(recvOrder.Menus); err != nil {
		p.RespError(c, nil, http.StatusUnprocessableEntity, "fail, Not Found Menu", err)
		return
	} else {
		resultOrder.Menus = menus
	}
	resultOrder.CustomerID = recvOrder.CustomerID
	if err = p.md.UpdateOrder(resultOrder, recvOrdernumber); err != nil {
		p.RespError(c, nil, http.StatusUnprocessableEntity, "fail update order", err)
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"res":      "ok",
			"contents": "The menu has changed..",
			"status":   "changed order",
			"Data":     resultOrder,
		})
	}
}

// GetSortedMenu godoc
// @Summary 메뉴 리스트 조회 및 정렬(추천/평점/주문수/최신)
// @Description 판매자가 추천설정을 한 메뉴들의 리스트를 가져올 수 있다.
// @Param	grade query string	false	"condition"
// @Param	createdat query string	false	"condition"
// @Param	1 query int	false	"orderby"
// @Success 200 {string} string "ok"
// @Failure 400 {string} string "err message"
// @Failure 422 {string} string "fail, Not Found Param"
// @Router /order [get]
func (p *Controller) GetSortedMenu(c *gin.Context) {
	query := model.QueryData{}
	if err := c.ShouldBindQuery(&query); err != nil {
		p.RespError(c, nil, http.StatusUnprocessableEntity, "fail, This is incorrect Query", err)
		return
	}
	if query.OrderBy != 1 {
		query.OrderBy = -1
	}
	result, err := p.md.GetSortedMenu(query)
	if err != nil {
		p.RespError(c, nil, http.StatusUnprocessableEntity, "fail, Not Found Param", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"res":  "ok",
		"data": result,
	})
}

// GetOrderHistory godoc
// @Summary 주문이 완료된 주문 내역을 보여줍니다.
// @Description 고객의 ID를 Parameter로 받아 주문 내역 중 주문이 완료된 주문 목록을 보여줍니다.
// @Param	customerID query string	false	"customerID"
// @Success 200 {string} string "ok"
// @Failure 400 {string} string "fail, Not Found Param"
// @Failure 422 {string} string "fail, not found customer"
// @Router /order/history/:customerID [get]
func (p *Controller) GetOrderHistory(c *gin.Context) {
	customerID := c.Param("customerID")
	if result, err := primitive.ObjectIDFromHex(customerID); err != nil {
		p.RespError(c, nil, http.StatusBadRequest, "fail, This is incorrect Param", err)
		c.Abort()
		return
	} else {
		if orders, err := p.md.GetOrdersInfoByUserID("his", result); err != nil {
			p.RespError(c, nil, http.StatusUnprocessableEntity, "fail, not found customer", err)
			return
		} else {
			if len(orders) == 0 {
				orders = nil
			}
			c.JSON(http.StatusOK, gin.H{
				"res":  "ok",
				"data": orders,
			})
		}
	}
}

// GetOrderInfo godoc
// @Summary 현재 고객이 주문한 주문 내용을 보여줍니다.
// @Description 고객의 ID를 Parameter로 받아 현재 고객이 주문한 주문 목록을 보여줍니다.
// @Param	customerID query string	false	"customerID"
// @Success 200 {string} string "ok"
// @Failure 400 {string} string "fail, Not Found Param"
// @Failure 422 {string} string "fail, not found customer"
// @Router /order/:customerID [get]
func (p *Controller) GetOrderInfo(c *gin.Context) {
	customerID := c.Param("customerID")
	if result, err := primitive.ObjectIDFromHex(customerID); err != nil {
		p.RespError(c, nil, http.StatusBadRequest, "fail, This is incorrect Param", err)
		c.Abort()
		return
	} else {
		if orders, err := p.md.GetOrdersInfoByUserID("", result); err != nil {
			p.RespError(c, nil, http.StatusUnprocessableEntity, "fail, not found customer", err)
			return
		} else {
			if len(orders) == 0 {
				c.JSON(http.StatusOK, gin.H{
					"res":  "ok",
					"data": nil,
				})
			}
			c.JSON(http.StatusOK, gin.H{
				"res":  "ok",
				"data": orders,
			})
		}
	}
}

// WriteReview godoc
// @Summary 리뷰를 작성할 수 있습니다.
// @Description 주문 번호를 확인 후 배송완료가 된 주문이면 리뷰와 평점을 작성할 수 있다.
// @Accept  json
// @Produce  json
// @Success 200 {string} string "ok"
// @Success 200 {string} string "Your order has not been completed."
// @Failure 400 {string} string "fail, This is incorrect JSON"
// @Failure 400 {string} string "You entered the wrong order number"
// @Failure 422 {string} string "fail, Not Found Order"
// @Failure 422 {string} string "fail, create review"
// @Router /history/review/:orderNumber [post]
func (p *Controller) WriteReview(c *gin.Context) {
	orderNumber := c.Param("orderNumber")
	recvReview := model.Review{}
	var err error
	if err = c.ShouldBindJSON(&recvReview); err != nil {
		p.RespError(c, nil, http.StatusBadRequest, "fail, This is incorrect JSON", err)
		return
	}

	var orderID primitive.ObjectID
	if orderID, err = primitive.ObjectIDFromHex(orderNumber); err != nil {
		p.RespError(c, nil, http.StatusBadRequest, "You entered the wrong order number.", err)
		return
	}

	if _, err := p.md.GetOrderByID(orderID); err != nil {
		p.RespError(c, nil, http.StatusUnprocessableEntity, "fail, Not Found Order", err)
		return
	} else {
		if orderStatus, err := p.md.GetOrderStatusByOrderID(orderID); err != nil {
			p.RespError(c, nil, http.StatusUnprocessableEntity, "fail, Not Found Order", err)
			return
		} else {
			if orderStatus != 7 {
				p.RespError(c, nil, http.StatusOK, "Your order has not been completed.", err)
				return
			}
			resultReview := model.Review{}
			resultReview.MenuId = recvReview.MenuId
			resultReview.Content = recvReview.Content
			resultReview.CustomerID = recvReview.CustomerID
			resultReview.Grade = recvReview.Grade
			resultReview.ID = primitive.NewObjectID()
			resultReview.CreatedAt = time.Now()
			resultReview.IsWrite = true
			if err = p.md.CreateReview(resultReview); err != nil {
				log.Println("fail insert review")
				p.RespError(c, nil, http.StatusUnprocessableEntity, "fail, create review", err)
				return
			} else {
				c.JSON(http.StatusOK, gin.H{
					"res": "ok",
				})
			}
		}
	}
}

// GetReview godoc
// @Summary 메뉴에 등록된 리뷰 목록을 확인할 수 있습니다.
// @Description Parameter로 menuID를 받아 해당 메뉴에 등록된 리뷰를 확인할 수 있습니다.
// @Accept  json
// @Produce  json
// @Success 200 {string} string "ok"
// @Success 400 {string} string "fail, Not Found Param."
// @Router /history/review/:orderNumber [post]
func (p *Controller) GetReview(c *gin.Context) {
	menuID := c.Param("menuID")
	if pMenuID, err := primitive.ObjectIDFromHex(menuID); err != nil {
		p.RespError(c, nil, http.StatusBadRequest, "fail, Not Found Param", err)
		c.Abort()
		return
	} else {
		if result, err := p.md.GetReviewByMenuID(pMenuID); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"res":  "ok",
				"data": nil,
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"res":  "ok",
				"data": result,
			})
		}
	}
}

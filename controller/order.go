package controller

import (
	"WBABEProject-04/model"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// 주문자

func (p *Controller) AddOrder(c *gin.Context) {
	// 메뉴추가
	// 주문한 order ID와
	// order의 post
	sOrderNumber := c.Param("ordernumber")
	var recvParam primitive.ObjectID
	var err error
	if recvParam, err = primitive.ObjectIDFromHex(sOrderNumber); err != nil {
		p.RespError(c, nil, 400, "fail, Not Found Param", nil)
		return
	}

	var recvOrder model.Order
	err = c.ShouldBindJSON(&recvOrder)
	if err != nil {
		p.RespError(c, nil, 400, "fail, Not Found Param", nil)
		c.Abort()
		return
	}
	tempOrder := model.Order{}
	if tempOrder, err = p.md.GetOrderByID(recvParam); err != nil {
		p.RespError(c, nil, 422, "fail, Not Found Order", err)
		return
	}

	if tempOrder.Status == model.Cooking || tempOrder.Status == model.InDelivery || tempOrder.Status == model.CompleteDelivery {
		newOrder := model.Order{}
		newOrder.CustomerID = recvOrder.CustomerID

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
		newOrder.Menus = menus
		newOrder.Status = model.Accepting
		newOrder.ID = primitive.NewObjectID()
		if err := p.md.CreateOrder(newOrder); err != nil {
			p.RespError(c, nil, 422, "fail, Not Found Param", err)
			return
		} else {
			c.JSON(http.StatusOK, gin.H{
				"res":  "ok",
				"Data": newOrder.ID,
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
			c.JSON(http.StatusOK, gin.H{
				"res":  "ok",
				"Data": recvParam,
			})
		}
	}
}

// OrderMenu godoc
// @Summary call OrderMenu, return ok by json.
// @Description 주문을 요청한다. 요청이 완료되면 주문 번호를 json 형태로 리턴한다.
// @Router /order [post]
func (p *Controller) OrderMenu(c *gin.Context) {
	customer := model.Customer{}
	err := c.ShouldBindJSON(&customer)
	if err != nil {
		p.RespError(c, nil, 400, "fail, Not Found Param", nil)
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
	menus := []model.Menu{}
	for _, order := range customer.Orders {
		tempOrder.Status = model.Accepting
		for _, menu := range order.Menus {
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
		tempOrder.Menus = menus
		tempOrder.Status = model.Receipt
		tempOrder.CreatedAt = time.Now()
		tempOrder.ID = primitive.NewObjectID()
		tempOrder.CustomerID = resultCustomer.ID
	}
	if err := p.md.CreateOrder(tempOrder); err != nil {
		p.RespError(c, nil, http.StatusUnprocessableEntity, "fail create order", err)
		return
	}

	resultOrders = append(resultOrders, tempOrder)
	resultCustomer.Orders = resultOrders
	if err := p.md.CreateCustomer(resultCustomer); err != nil {
		p.RespError(c, nil, http.StatusUnprocessableEntity, "fail create order", err)
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

// GetSortedMenu godoc
// @Summary 메뉴 리스트 조회 및 정렬(추천/평점/주문수/최신)
// @Description 판매자가 추천설정을 한 메뉴들의 리스트를 가져올 수 있다.
// @Param	grade query string	false	"condition"
// @Param	favorite query string	false	"condition"
// @Param	count query string	false	"condition"
// @Param	createdat query string	false	"condition"
// @Param	1 query int	false	"orderby"
// @Success 200 {string} string "ok"
// @Failure 400 {string} string "err message"
// @Failure 422 {string} string "fail, Not Found Param"
// @Router /order [get]
func (p *Controller) GetSortedMenu(c *gin.Context) {
	query := model.QueryData{}

	if err := c.ShouldBindQuery(&query); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	if query.OrderBy != 1 {
		query.OrderBy = -1
	}
	result, err := p.md.GetSortedMenu(query)
	if err != nil {
		p.RespError(c, nil, http.StatusUnprocessableEntity, "fail, Not Found Param", err)
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"res":  "ok",
		"data": result,
	})
}

func (p *Controller) GetOrderHistory(c *gin.Context) {
	customerID := c.Param("customerID")
	if result, err := primitive.ObjectIDFromHex(customerID); err != nil {
		p.RespError(c, nil, 400, "fail, Not Found Param", err)
		c.Abort()
		return
	} else {
		if orders, err := p.md.GetOrdersInfoByUserID("his", result); err != nil {
			p.RespError(c, nil, 400, "fail, not found customer", err)
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

func (p *Controller) GetOrderInfo(c *gin.Context) {
	customerID := c.Param("customerID")
	if result, err := primitive.ObjectIDFromHex(customerID); err != nil {
		p.RespError(c, nil, 400, "fail, Not Found Param", err)
		c.Abort()
		return
	} else {
		if orders, err := p.md.GetOrdersInfoByUserID("", result); err != nil {
			p.RespError(c, nil, 400, "fail, not found customer", err)
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
// @Summary call WriteReview, return ok by json.
// @Description 주문 번호를 확인 후 배송완료가 된 주문이면 리뷰와 평점을 작성할 수 있다.
// @name WriteReview
// @Accept  json
// @Produce  json
// @Router /history/review/:orderNumber [post]
func (p *Controller) WriteReview(c *gin.Context) {
	orderNumber := c.Param("orderNumber")
	recvReview := model.Review{}
	var err error
	if err = c.ShouldBindJSON(&recvReview); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	var orderID primitive.ObjectID
	if orderID, err = primitive.ObjectIDFromHex(orderNumber); err != nil {
		p.RespError(c, nil, 400, "You entered the wrong order number.", err)
		return
	}

	if _, err := p.md.GetOrderByID(orderID); err != nil {
		p.RespError(c, nil, 400, "fail, Not Found Order", err)
		return
	} else {
		if orderStatus, err := p.md.GetOrderStatusByOrderID(orderID); err != nil {
			p.RespError(c, nil, 400, "fail, Not Found Order", err)
			return
		} else {
			if orderStatus != 7 {
				p.RespError(c, nil, 200, "Your order has not been completed.", err)
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
				p.RespError(c, nil, 422, "fail, create review", err)
				return
			} else {
				c.JSON(http.StatusOK, gin.H{
					"res": "ok",
				})
			}
		}
	}
}

func (p *Controller) GetReview(c *gin.Context) {
	menuID := c.Param("menuID")
	if pMenuID, err := primitive.ObjectIDFromHex(menuID); err != nil {
		p.RespError(c, nil, 400, "fail, Not Found Param", err)
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

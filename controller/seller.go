package controller

import (
	"WBABEProject-04/logger"
	"WBABEProject-04/model"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 판매자

// RegisterMenu godoc
// @Summary call RegisterMenu, return ok by json.
// @Description 메뉴를 등록한다.
// @Accept  json
// @Produce  json
// @Router /menu [post]
func (p *Controller) RegisterMenu(c *gin.Context) {
	// 필수 정보
	// 이름(name)
	// 가격(price)
	// 수량(quantity)
	// 원산지(Origin)
	fmt.Println("RegisterMenu")
	pMenu := model.Menu{}

	if err := c.ShouldBindJSON(&pMenu); err != nil {
		p.RespError(c, nil, 400, "fail, Not Found Param", nil)
		c.Abort()
		return
	}

	if len(pMenu.Name) <= 0 {
		p.RespError(c, nil, 400, "fail, Not Found Param", nil)
		c.Abort()
		return
	}
	if len(pMenu.Origin) <= 0 {
		pMenu.Origin = "국내산"
	}

	// pReview := []model.Review{}
	pMenu.Review = nil

	_, err := p.md.GetOneMenu(pMenu.Name)

	// 이미 등록된 메뉴가 있으면
	if err == nil {
		p.RespError(c, nil, http.StatusUnprocessableEntity, "already resistery menu", nil)
		return
	}

	err = p.md.CreateMenu(pMenu)
	if err != nil {
		p.RespError(c, nil, http.StatusUnprocessableEntity, "parameter not found", err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"result": "ok",
	})
	c.Next()
}

// DelMenu godoc
// @Summary call DelMenu, return ok by json.
// @Description 메뉴의 이름을 파라미터로 받아 해당 메뉴를 삭제하는 기능
// @Router /order/menu/:name [delete]
func (p *Controller) DelMenu(c *gin.Context) {
	logger.Debug("DelMenu")
	sMenu := c.Param("menu")
	if len(sMenu) <= 0 {
		p.RespError(c, nil, http.StatusUnprocessableEntity, "parameter not found", nil)
		return
	}

	_, err := p.md.GetOneMenu(sMenu)
	if err != nil {
		p.RespError(c, nil, http.StatusUnprocessableEntity, "exist resistery menu", nil)
		return
	}

	if err := p.md.DeleteMenu(sMenu); err != nil {
		p.RespError(c, nil, http.StatusUnprocessableEntity, "fail delete db", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": "ok",
	})
	c.Next()
}

// // UpdateMenu godoc
// // @Summary call UpdateMenu, return ok by json.
// // @Description 기존 메뉴의 정보를 변경할 수 있다.
// // @Router /menu [post]
func (p *Controller) UpdateMenu(c *gin.Context) {
	var recvMenu model.Menu
	err := c.ShouldBindJSON(&recvMenu)
	if err != nil {
		p.RespError(c, nil, 400, "fail, Not Found Param", nil)
		c.Abort()
		return
	}

	menu, err := p.md.GetOneMenu(recvMenu.Name)
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

	// 가져온 주문 id 값을 통해 상태 값 변경 (접수 이상으로 진행된 것)

	// orderList, err := p.md.UpdateStatus()
	// if err != nil {
	// 	return
	// }
	// fmt.Println(orderList)
}

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

// // RegisterStore godoc
// // @Summary call RegisterStore, return ok by json.
// // @Description 가게를 등록하는 기능
// // @Accept  json
// // @Produce  json
// // @Router /order/seller/store [post]
// func (p *Controller) RegisterStore(c *gin.Context) {
// 	store := model.Store{}
// 	fmt.Println(store)
// 	if err := c.ShouldBindJSON(&store); err != nil {
// 		panic(err)
// 	}
// 	fmt.Println("store.Name len : ", len(store.Name))
// 	if len(store.Name) <= 0 {
// 		p.RespError(c, nil, http.StatusUnprocessableEntity, "parameter not found", nil)
// 		return
// 	}
// 	fmt.Println(store.Status)
// 	fmt.Println(store.Menus)

// 	if len(store.Menus) != 0 || len(store.Reviews) != 0 {
// 		p.RespError(c, nil, http.StatusUnprocessableEntity, "parameter not found", nil)
// 		return
// 	}
// 	fmt.Println("store.Status len : ", len(store.Status))
// 	if len(store.Status) > 0 {
// 		p.RespError(c, nil, http.StatusUnprocessableEntity, "parameter not found", nil)
// 		return
// 	}

// 	p.md.CreateStore(store)
// }

// func (p *Controller) GetStoreInfo(c *gin.Context) {
// 	// filter := bson.M{"name":}
// 	sName := c.Param("name")
// 	if per, err := p.md.GetOneStore(sName); err != nil {
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

// // DelMenu godoc
// // @Summary call DelMenu, return ok by json.
// // @Description 메뉴의 이름을 파라미터로 받아 해당 메뉴를 삭제하는 기능
// // @Router /order/menu/:name [delete]
// func (p *Controller) DelMenu(c *gin.Context) {
// 	logger.Debug("DelMenu")
// 	smenu := c.Param("menu")
// 	if len(smenu) <= 0 {
// 		p.RespError(c, nil, http.StatusUnprocessableEntity, "parameter not found", nil)
// 		return
// 	}

// 	_, err := p.md.GetOneMenu("mune", smenu)
// 	if err != nil {
// 		p.RespError(c, nil, http.StatusUnprocessableEntity, "exist resistery person", nil)
// 		return
// 	}

// 	if err := p.md.DeleteMenu(smenu); err != nil {
// 		p.RespError(c, nil, http.StatusUnprocessableEntity, "fail delete db", err)
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"result": "ok",
// 	})
// 	c.Next()
// }
// func NewController(rep *model.Model) (*Controller, error) {

// 	r := &Controller{md: rep}
// 	return r, nil
// }
// func (p *Controller) RegisterMenu(c *gin.Context) {
// 	logger.Debug("RegisterMenu")
// 	name := c.PostForm("name")         // (필수)
// 	order := c.PostForm("order")       // 주문 가능 여부 (Default: 불가능)
// 	quantity := c.PostForm("quantity") // 주문가능 개수(Default: infinity)
// 	origin := c.PostForm("origin")     // 원산지 (필수, Default: 국내산)
// 	price := c.PostForm("price")       // 가격 (필수)
// 	spicy := c.PostForm("spicy")       // 맵기 (Default: normal)

// 	var nPrice uint
// 	tempPrice, err := strconv.Atoi(price)
// 	if tempPrice < 0 {
// 		nPrice = 0
// 	}
// 	if err != nil {
// 		nPrice = 0
// 	}

// 	nQuantity, err := strconv.Atoi(quantity)
// 	if err != nil {
// 		// 수량을 따로 정해두지 않음
// 		nQuantity = 1000000
// 	}

// 	// 필수 정보로 들어가야할 정보가 없으면
// 	if len(name) <= 0 || len(origin) <= 0 {
// 		p.RespError(c, nil, http.StatusUnprocessableEntity, "parameter not found")
// 	}

// 	if len(origin) <= 0 {
// 		origin = "국내산"
// 	}

// 	// 주문 불가능 여부 : (Default: 가능(false))
// 	bOrder, err := strconv.ParseBool(order)
// 	if err != nil {
// 		bOrder = false
// 	}
// 	if len(spicy) <= 0 {
// 		spicy = "Normal"
// 	}

// 	menu, _ := p.md.GetOneMenu("name", name)

// 	// 이미 등록된 메뉴가 있으면
// 	if menu != (model.Menu{}) {
// 		p.RespError(c, nil, http.StatusUnprocessableEntity, "already resistery menu", nil)
// 		return
// 	}

// 	// 입력된 정보로 메뉴 생성
// 	req := model.Menu{Name: name, Order: bOrder, Quantity: nQuantity, Spicy: spicy, Origin: origin, Price: nPrice}

// 	if err := p.md.CreateMenu(req); err != nil {
// 		p.RespError(c, nil, http.StatusUnprocessableEntity, "parameter not found", err)
// 		return
// 	}
// 	// 요청 성공 시
// 	c.JSON(http.StatusCreated, gin.H{
// 		"result": "ok",
// 	})
// 	c.Next()
// }

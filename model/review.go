package model

type Review struct {
	Content  string   `json:"content" bson:"content"`
	Menus    Menu     `json:"menus" bson:"menus"`
	Customer Customer `json:"customer" bson:"customer"`
}

// // 메뉴이름을 받아 메뉴를 가져온다.
// func (m *Model) GetOneMenu(flag, elem string) (Menu, error) {

// 	logger.Debug("seller > GetOneMenu")
// 	opts := []*options.FindOneOptions{}
// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()

// 	var filter bson.M

// 	if flag == "name" {
// 		filter = bson.M{"name": elem}
// 	}
// 	var menus Menu
// 	if err := m.collectionSeller.FindOne(ctx, filter, opts...).Decode(&menus); err != nil {
// 		return menus, err
// 	} else {
// 		return menus, nil
// 	}
// }

// // 메뉴를 생성한다.
// func (m *Model) CreateMenu(menus Menu) error {
// 	logger.Debug("seller > CreateMenu")
// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()

// 	if _, err := m.collectionSeller.InsertOne(ctx, menus); err != nil {
// 		log.Println("fail insert new menu")
// 		return fmt.Errorf("fail, insert")
// 	}
// 	return nil
// }

// func (m *Model) DeleteMenu(smenu string) error {
// 	logger.Debug("seller > DeleteMenu")
// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()

// 	filter := bson.M{"name": smenu}
// 	if res, err := m.collectionSeller.DeleteOne(ctx, filter); res.DeletedCount <= 0 {
// 		return fmt.Errorf("could not delete, not found menu %s", smenu)
// 	} else if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (m *Model) UpdateMenu(menu Menu) error {
// 	fmt.Println("UpdateMenu : ", menu)
// 	filter := bson.M{"name": menu.Name}
// 	update := bson.M{
// 		"$set": bson.M{
// 			"order":    menu.Order,
// 			"quantity": menu.Quantity,
// 			"price":    menu.Price,
// 			"spicy":    menu.Spicy,
// 			"origin":   menu.Origin,
// 		},
// 	}
// 	if _, err := m.collectionSeller.UpdateOne(context.Background(), filter, update); err != nil {
// 		return err
// 	}
// 	fmt.Println(">>??")
// 	return nil

// }

// func (m *Model) GetMenuList() []Menu {

// 	logger.Debug("seller > GetMenuList")
// 	fmt.Println("GetMenuList")
// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()

// 	cursor, err := m.collectionSeller.Find(ctx, bson.M{})
// 	if err != nil {
// 		panic(err)
// 	}

// 	var menus []Menu
// 	if err = cursor.All(ctx, &menus); err != nil {
// 		panic(err)
// 	}

// 	if err != nil {
// 		panic(err)
// 	}

// 	return menus
// }

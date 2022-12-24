package model

type Spiciness int

const (
	Normal = iota + 1
	Spicy
	Hot
	Hell
	Burning
)

type Menu struct {
	Name      string    `json:"name" bson:"name"`           // 메뉴 이름
	Available bool      `json:"available" bson:"available"` // 주문 가능 여부
	Quantity  int       `json:"quantity" bson:"quantity"`   // 수량
	Grade     int       `json:"grade" bson:"grade"`
	Origin    string    `json:"origin" bson:"origin"`       // 원산지
	Price     uint      `json:"price" bson:"price"`         // 가격
	Spiciness Spiciness `json:"spiciness" bson:"spiciness"` // 맵기 정도
	Referrals bool      `json:"referrals" bson:"referrals"` // 추천 여부
	Review    []Review  `json:"review" bson:"review"`
}

package controller

import "WBABEProject-04/model"

type Controller struct {
	md *model.Model
}

func NewController(rep *model.Model) (*Controller, error) {
	r := &Controller{md: rep}
	return r, nil
}

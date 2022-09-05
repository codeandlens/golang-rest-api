package requestmodel

import (
	"example/go-demo-api/model"
)

type Requestmodel struct {
	Code string       `json:"code"`
	Desc string       `json:"desc"`
	Data []model.User `json:"data"`
}

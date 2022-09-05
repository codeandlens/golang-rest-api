package responsemodel

import "example/go-demo-api/model"

type Responsemodel struct {
	Code string       `json:"code"`
	Desc string       `json:"desc"`
	Data []model.User `json:"data"`
}

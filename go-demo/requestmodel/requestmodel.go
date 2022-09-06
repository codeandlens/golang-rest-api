package requestmodel

import (
	"example/go-demo-api/model"
)

type Requestmodel struct {
	Data []model.User `json:"data"`
}

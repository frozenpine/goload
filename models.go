package main

import (
	"encoding/json"
	"gitlab.quantdo.cn/yuanyang/goload/utils"
)

// Order order define for new
type Order struct {
	Symbol   string          `csv:"symbol" json:"symbol"`
	Price    float64         `csv:"price" json:"price"`
	Quantity int64           `csv:"quantity" json:"quantity"`
	Side     utils.OrderSide `csv:"side" json:"side"`
}

// GetQuantity get normalized quantity
func (ord *Order) GetQuantity() int64 {
	if ord.Quantity < 0 {
		return ord.Quantity
	}

	return ord.Quantity * ord.Side.Value()
}

func (ord *Order) String() string {
	result, _ := json.Marshal(*ord)

	return string(result)
}

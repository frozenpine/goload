package main

import (
	"gitlab.quantdo.cn/yuanyang/goload/utils"
)

// Order order define for new
type Order struct {
	Symbol   string          `csv:"symbol"`
	Price    float64         `csv:"price"`
	Quantity int64           `csv:"quantity"`
	Side     utils.OrderSide `csv:"side"`
}

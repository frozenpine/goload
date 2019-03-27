package main

import (
	"gitlab.quantdo.cn/yuanyang/goload/utils"
)

const (
	defaultHost        = "http://trade"
	defaultURI         = "/api/v1"
	defaultSymbol      = "XBTUSD"
	defaultQuantity    = int64(1)
	defaultMaxQuantity = int64(100)
	defaultPrice       = float64(3536)
	defaultPrecision   = 2
	defaultSide        = utils.Buy

	defaultIdentity = "yuanyang@quantdo.com.cn"
	defaultPassword = "quantdo123456"
)

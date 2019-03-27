package main

import (
	"encoding/json"
	"testing"

	"gitlab.quantdo.cn/yuanyang/goload/utils"
)

func TestOrderMarshal(t *testing.T) {
	ord := Order{
		Symbol:   "XBTUSD",
		Price:    3536.0,
		Quantity: 1,
		Side:     utils.Buy,
	}

	result, err := json.Marshal(ord)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(string(result))
	}
}

func TestOrderUnmarshal(t *testing.T) {
	ordString := `{"symbol":"XBTUSD","price":3721,"quantity":100,"side":"Sell"}`

	var ord Order

	if err := json.Unmarshal([]byte(ordString), &ord); err != nil {
		t.Error(err)
	} else if ord.Symbol != "XBTUSD" || ord.Price != 3721.0 || ord.Quantity != 100 || ord.Side != utils.Sell {
		t.Error("miss-match data unmarshaled.")
	}
}

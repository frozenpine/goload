package utils

import (
	"bytes"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
)

// OrderSide order side
type OrderSide string

const (
	// Buy buy side
	Buy OrderSide = "Buy"
	// Sell sell side
	Sell OrderSide = "Sell"
)

// String order side string
func (s OrderSide) String() string {
	return string(s)
}

// Value order side value
func (s OrderSide) Value() int64 {
	switch s {
	case Buy:
		return 1
	case Sell:
		return -1
	default:
		panic(ErrSide)
	}
}

// Opposite get opposite order side
func (s OrderSide) Opposite() OrderSide {
	switch s {
	case Buy:
		return Sell
	case Sell:
		return Buy
	default:
		panic(ErrSide)
	}
}

// UnmarshalCSV unmarshal csv column to OrderSide
func (s *OrderSide) UnmarshalCSV(value string) error {
	return s.Set(value)
}

// MarshalCSV marshal to csv column
func (s *OrderSide) MarshalCSV() string {
	return (*s).String()
}

// UnmarshalJSON unmarshal from json string
func (s *OrderSide) UnmarshalJSON(data []byte) error {
	return s.Set(strings.Trim(string(data), "\""))
}

// MarshalJSON marshal to json string
func (s *OrderSide) MarshalJSON() ([]byte, error) {
	var buff bytes.Buffer
	buff.WriteString((*s).String())

	return buff.Bytes(), nil
}

// Set set value for flag
func (s *OrderSide) Set(value string) error {
	switch value {
	case "Buy", "buy":
		*s = Buy
		return nil
	case "Sell", "sell":
		*s = Sell
		return nil
	default:
		return ErrSide
	}
}

// CheckSymbol validate order symbol
func CheckSymbol(symbol string) error {
	if symbol == "" {
		return ErrSymbol
	}

	return nil
}

// CheckPrice validate order price
func CheckPrice(price float64) error {
	if price <= 0 {
		return ErrPrice
	}

	return nil
}

// CheckQuantity validate order quantity
func CheckQuantity(qty int64) error {
	if qty == 0 {
		return ErrQuantity
	}

	return nil
}

// MatchSide match side with quantity
func MatchSide(side *OrderSide, qty int64) error {
	switch *side {
	case "Buy", "buy":
		if qty < 0 {
			return ErrMissMatchQtySide
		}
		*side = Buy
	case "Sell", "sell":
		if qty > 0 {
			return ErrMissMatchQtySide
		}
		*side = Sell
	case "":
		if qty > 0 {
			*side = Buy
		} else {
			*side = Sell
		}
	default:
		return ErrSide
	}

	return nil
}

// RandomSide generate random side
func RandomSide() OrderSide {
	choice := rand.Intn(10)
	if choice < 5 {
		return Buy
	}

	return Sell
}

// RandomPrice generate random price on basePrice
func RandomPrice(price *float64, prec int, basePrice float64) error {
	if prec <= 0 {
		prec = 2
	}

	out, err := strconv.ParseFloat(
		fmt.Sprintf("%."+strconv.Itoa(prec)+"f", basePrice+rand.Float64()), 64)
	if err != nil {
		return err
	}

	*price = out

	return nil
}

// RandomQuantity generate random quantity (0, max]
func RandomQuantity(quantity *int64, max int64) error {
	if max <= 1 {
		*quantity = 1
	} else {
		*quantity = rand.Int63n(max) + 1
	}

	return nil
}

package utils

import (
	"fmt"
	"math/rand"
	"strconv"
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
func (s *OrderSide) String() string {
	return string(*s)
}

// Value order side value
func (s *OrderSide) Value() int64 {
	switch *s {
	case Buy:
		return 1
	case Sell:
		return -1
	default:
		return 1
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
func MatchSide(side *string, qty int64) error {
	switch *side {
	case "Buy", "buy":
		if qty < 0 {
			return ErrMissMatchQtySide
		}
		*side = Buy.String()
	case "Sell", "sell":
		if qty > 0 {
			return ErrMissMatchQtySide
		}
		*side = Sell.String()
	case "":
		if qty > 0 {
			*side = Buy.String()
		} else {
			*side = Sell.String()
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

// RandomPrice generate random price
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

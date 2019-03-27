package utils

import (
	"errors"
)

var (
	// ErrSymbol invalid symbol
	ErrSymbol = errors.New("symbol can't be empty")

	// ErrPrice invalid price
	ErrPrice = errors.New("order price should be positive")

	// ErrQuantity invalid quantity
	ErrQuantity = errors.New("order quantity can't be ZERO")

	// ErrMissMatchQtySide quantity miss-match with side
	ErrMissMatchQtySide = errors.New("order quantity miss-match with side")

	// ErrSide invalid side
	ErrSide = errors.New("side is either \"Buy\" or \"Sell\"")

	// ErrIdentity invalid login identity
	ErrIdentity = errors.New("identity should either be email or mobile")
)

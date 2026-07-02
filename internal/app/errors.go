package app

import "errors"

// ErrInsufficientCustodyBalance is returned when a return exceeds the current custody quantity.
var ErrInsufficientCustodyBalance = errors.New("insufficient custody balance")

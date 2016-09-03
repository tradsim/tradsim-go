package models

import "fmt"

// OrderStatus defines the status of the order
type OrderStatus uint8

// The various statuses of a order
const (
	Pending OrderStatus = iota
	PartiallyFilled
	FullyFilled
	OverFilled
	Cancelled
)

// Order status string
const (
	PendingText         = "Pending"
	PartiallyFilledText = "PartiallyFilled"
	FullyFilledText     = "FullyFilled"
	OverFilledText      = "OverFilled"
	CancelledText       = "Cancelled"
)

func (o OrderStatus) String() string {
	switch o {
	case Pending:
		return PendingText
	case PartiallyFilled:
		return PartiallyFilledText
	case FullyFilled:
		return FullyFilledText
	case OverFilled:
		return OverFilledText
	case Cancelled:
		return CancelledText
	default:
		return fmt.Sprintf("Not mapped value %d", o)
	}
}

// IsTradeable return true if status allows trading
func (o OrderStatus) IsTradeable() bool {
	switch o {
	case Pending, PartiallyFilled:
		return true
	default:
		return false
	}
}

// OrderStatusFromString returns a order status from string
func OrderStatusFromString(value string) (OrderStatus, error) {
	switch value {
	case PendingText:
		return Pending, nil
	case PartiallyFilledText:
		return PartiallyFilled, nil
	case FullyFilledText:
		return FullyFilled, nil
	case OverFilledText:
		return OverFilled, nil
	case CancelledText:
		return Cancelled, nil
	default:
		return 9, fmt.Errorf("Not mapped %s", value)
	}
}

// ResolveStatus resolves the status based on the quantities
func ResolveStatus(q uint, t uint) OrderStatus {
	if t == uint(0) {
		return Pending
	} else if t < q {
		return PartiallyFilled
	} else if t == q {
		return FullyFilled
	}
	return OverFilled
}

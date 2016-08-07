package models

import (
	"fmt"
)

// TradeDirection defines buy or sell
type TradeDirection uint8

// Buy or sell
const (
	Buy TradeDirection = iota
	Sell
)

// Buy and sell string
const (
	BuyText  = "Buy"
	SellText = "Sell"
)

func (t TradeDirection) String() string {
	switch t {
	case Buy:
		return BuyText
	case Sell:
		return SellText
	default:
		return fmt.Sprintf("Not mapped value %d", t)
	}
}

// TradeDirectionFromString returns a trade direction from string
func TradeDirectionFromString(value string) (TradeDirection, error) {
	switch value {
	case BuyText:
		return Buy, nil
	case SellText:
		return Sell, nil
	default:
		return 9, fmt.Errorf("Not mapped %s", value)
	}
}

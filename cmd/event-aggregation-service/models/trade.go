package models

import (
	"time"

	"github.com/satori/go.uuid"
)

// Trade defines a trade
type Trade struct {
	ID       int64
	OrderID  uuid.UUID
	Price    float64
	Quantity uint
	Occured  time.Time
}

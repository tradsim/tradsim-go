package models

import "time"

// Position defines the position
type Position struct {
	ID       int64
	Symbol   string
	Quantity int
	Updated  time.Time
}

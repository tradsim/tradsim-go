package models

import "time"

// Position defines the position
type Position struct {
	Symbol   string
	Quantity int
	Updated  time.Time
}

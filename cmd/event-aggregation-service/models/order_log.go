package models

import (
	"time"

	"github.com/satori/go.uuid"
)

// OrderLog defines a order log
type OrderLog struct {
	ID      int64
	OrderID uuid.UUID
	Action  string
	Occured time.Time
}

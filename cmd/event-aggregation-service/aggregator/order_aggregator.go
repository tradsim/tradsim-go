package aggregator

import (
	"time"

	"github.com/tradsim/tradsim-go/cmd/event-aggregation-service/models"
	commonmodels "github.com/tradsim/tradsim-go/models"
)

// OrderAggregator interface
type OrderAggregator interface {
	Aggregate(s string, orders []models.Order) models.Position
}

// OrderAggregatorImpl aggregates order into a position
type OrderAggregatorImpl struct {
}

// NewOrderAggregator returns a new order aggregator
func NewOrderAggregator() *OrderAggregatorImpl {
	return &OrderAggregatorImpl{}
}

// Aggregate order into a position
func (oa OrderAggregatorImpl) Aggregate(s string, orders []models.Order) models.Position {
	q := 0
	var t time.Time

	for _, o := range orders {
		if o.Symbol != s {
			continue
		}
		if o.Direction == commonmodels.Buy {
			q += int(o.Traded)
		} else {
			q -= int(o.Traded)
		}

		if t.IsZero() {
			t = o.Updated
		} else {
			if o.Updated.After(t) {
				t = o.Updated
			}
		}
	}

	return models.Position{s, q, t}
}

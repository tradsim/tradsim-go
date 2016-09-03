package aggregator

import (
	"fmt"
	"log"

	incmodel "github.com/mantzas/incata/model"
	uuid "github.com/satori/go.uuid"
	"github.com/tradsim/tradsim-go/events"
	"github.com/tradsim/tradsim-go/models"
)

// Aggregator interface
type Aggregator interface {
	Aggregate(events []incmodel.Event) (models.Order, error)
}

// EventAggregator aggregates events to order
type EventAggregator struct {
}

// NewAggregator creates a new aggregator
func NewAggregator() *EventAggregator {
	return &EventAggregator{}
}

// Aggregate events to a order
func (ea EventAggregator) Aggregate(evs []incmodel.Event) (models.Order, error) {
	var or models.Order

	for _, ev := range evs {

		switch ev.EventType {
		case string(events.OrderAcceptedType):
			err := aggregateAccepted(&or, ev)
			if err != nil {
				return models.Order{}, err
			}
		case string(events.OrderAmendedType):
			err := aggregateAmended(&or, ev)
			if err != nil {
				return models.Order{}, err
			}
		case string(events.OrderCancelledType):
			err := aggregateCanceled(&or, ev)
			if err != nil {
				return models.Order{}, err
			}
		case string(events.OrderTradedType):
			err := aggregateTraded(&or, ev)
			if err != nil {
				return models.Order{}, err
			}
		default:
			return models.Order{}, fmt.Errorf("Event type not supported %s", ev.EventType)
		}
	}
	return or, nil
}

func aggregateAccepted(o *models.Order, ev incmodel.Event) error {

	accepted, ok := ev.Payload.(events.OrderAccepted)
	if !ok {
		return fmt.Errorf("type assertion to %s failed", ev.EventType)
	}

	dir, err := models.TradeDirectionFromString(accepted.Direction)
	if err != nil {
		return err
	}

	orderID, err := uuid.FromString(accepted.OrderID)
	if err != nil {
		return err
	}

	o.ID = orderID
	o.Symbol = accepted.Symbol
	o.Price = accepted.Price
	o.Quantity = accepted.Quantity
	o.Traded = 0
	o.Direction = dir
	o.Status = models.Pending
	log.Printf("Accepted aggregation succeeded. %s", o.String())
	return nil
}

func aggregateAmended(o *models.Order, ev incmodel.Event) error {
	amended, ok := ev.Payload.(events.OrderAmended)
	if !ok {
		return fmt.Errorf("type assertion to %s failed", ev.EventType)
	}
	o.Quantity = amended.Quantity
	o.UpdateStatus()
	log.Printf("Amended aggregation succeeded. %s", o.String())
	return nil
}

func aggregateCanceled(o *models.Order, ev incmodel.Event) error {
	_, ok := ev.Payload.(events.OrderCancelled)
	if !ok {
		return fmt.Errorf("type assertion to %s failed", ev.EventType)
	}
	o.Status = models.Cancelled
	log.Printf("Canceled aggregation succeeded. %s", o.String())
	return nil
}

func aggregateTraded(o *models.Order, ev incmodel.Event) error {
	traded, ok := ev.Payload.(events.OrderTraded)
	if !ok {
		return fmt.Errorf("type assertion to %s failed", ev.EventType)
	}
	o.Trade(traded.Quantity)
	log.Printf("Traded aggregation succeeded. %s", o.String())
	return nil
}

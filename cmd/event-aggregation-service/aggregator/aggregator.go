package aggregator

import (
	"fmt"
	"log"

	incmodel "github.com/mantzas/incata/model"
	uuid "github.com/satori/go.uuid"
	"github.com/tradsim/tradsim-go/cmd/event-aggregation-service/models"
	"github.com/tradsim/tradsim-go/events"
	commonmodels "github.com/tradsim/tradsim-go/models"
)

// Aggregator interface
type Aggregator interface {
	Aggregate(events []incmodel.Event) (models.Order, error)
}

// EventAggregator aggregates events to order
type EventAggregator struct {
}

// NewEventAggregator creates a new event aggregator
func NewEventAggregator() *EventAggregator {
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

	dir, err := commonmodels.TradeDirectionFromString(accepted.Direction)
	if err != nil {
		return err
	}

	orderID, err := uuid.FromString(accepted.OrderID)
	if err != nil {
		return err
	}

	*o = *models.NewOrder(orderID, accepted.Symbol, accepted.Price, accepted.Quantity, dir, commonmodels.Pending, accepted.Occured)

	log.Print("Accepted aggregation succeeded")
	return nil
}

func aggregateAmended(o *models.Order, ev incmodel.Event) error {
	amended, ok := ev.Payload.(events.OrderAmended)
	if !ok {
		return fmt.Errorf("type assertion to %s failed", ev.EventType)
	}
	o.Amend(amended.Quantity, amended.Occured)
	log.Print("Amended aggregation succeeded")
	return nil
}

func aggregateCanceled(o *models.Order, ev incmodel.Event) error {
	cancelled, ok := ev.Payload.(events.OrderCancelled)
	if !ok {
		return fmt.Errorf("type assertion to %s failed", ev.EventType)
	}
	o.Cancel(cancelled.Occured)
	log.Print("Canceled aggregation succeeded")
	return nil
}

func aggregateTraded(o *models.Order, ev incmodel.Event) error {
	traded, ok := ev.Payload.(events.OrderTraded)
	if !ok {
		return fmt.Errorf("type assertion to %s failed", ev.EventType)
	}
	o.Trade(traded.Price, traded.Quantity, traded.Occured)
	log.Print("Traded aggregation succeeded")
	return nil
}

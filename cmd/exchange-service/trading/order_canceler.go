package trading

import (
	"log"
	"sync"
	"time"

	"github.com/satori/go.uuid"
	"github.com/tradsim/tradsim-go/events"
	"github.com/tradsim/tradsim-go/models"
)

// Canceller interface
type Canceller interface {
	Cancel(book *models.OrderBook, orderID uuid.UUID) bool
}

// OrderCanceller for canceling orders
type OrderCanceller struct {
	publisher events.EventPublisher
	mu        sync.Mutex
}

// NewOrderCanceller creates a order canceller
func NewOrderCanceller(publisher events.EventPublisher) *OrderCanceller {
	return &OrderCanceller{publisher, sync.Mutex{}}
}

// Cancel order by id
func (oc *OrderCanceller) Cancel(book *models.OrderBook, orderID uuid.UUID) bool {
	oc.mu.Lock()
	defer oc.mu.Unlock()

	order, ok := book.Orders[orderID]
	if !ok {
		log.Printf("Order with id %s not found", orderID)
		return false
	}

	if !order.Status.IsTradeable() {
		return false
	}

	order.Status = models.Cancelled
	oc.publishCancelledEvent(orderID)
	return true
}

func (oc *OrderCanceller) publishCancelledEvent(ID uuid.UUID) {

	ev := events.NewOrderCancelled(ID.String(), time.Now().UTC(), uint(1))
	env, err := events.NewOrderEventEnvelope(ev, ev.EventType)

	if err != nil {
		log.Printf("Failed to create envelope: %s", err.Error())
	}

	err = oc.publisher.Publish(env)
	if err != nil {
		log.Printf("Failed to publish cancelled event: %s", ev.String())
	}
}

package trading

import (
	"errors"
	"log"
	"sync"
	"time"

	"github.com/satori/go.uuid"
	"github.com/tradsim/tradsim-go/events"
	"github.com/tradsim/tradsim-go/models"
)

// Amender interface
type Amender interface {
	Amend(book *models.OrderBook, order *models.Order) bool
}

// OrderAmender amends a order in the book
type OrderAmender struct {
	publisher events.EventPublisher
	mu        sync.Mutex
}

// NewOrderAmender creates a new order amender
func NewOrderAmender(publisher events.EventPublisher) *OrderAmender {
	return &OrderAmender{publisher, sync.Mutex{}}
}

// Amend a order in the order book
func (oa *OrderAmender) Amend(book *models.OrderBook, order *models.Order) bool {

	oa.mu.Lock()
	defer oa.mu.Unlock()

	prices, ok := book.Symbols[order.Symbol]

	if !ok {
		log.Printf("Symbol %s not found", order.Symbol)
		return false
	}

	found, i := findPrice(prices, order.Price)

	if !found {
		log.Printf("Price %f not found", order.Price)
		return false
	}

	var amended = false

	if order.Direction == models.Buy {

		for _, o := range prices[i].Buy.Orders {

			if o.ID == order.ID {

				d, err := oa.getAmendQuantity(o.Quantity, order.Quantity)
				if err != nil {
					log.Printf("%s", err)
					break
				}

				prices[i].Buy.Quantity += d
				o.Amend(d)
				amended = true
				break
			}
		}

	} else {

		for _, o := range prices[i].Sell.Orders {

			if o.ID == order.ID {

				d, err := oa.getAmendQuantity(o.Quantity, order.Quantity)
				if err != nil {
					log.Printf("%s", err)
					break
				}

				prices[i].Sell.Quantity += d
				o.Amend(d)
				amended = true
				break
			}
		}
	}

	if amended {
		oa.publishAmendEvent(order.ID, order.Quantity)
	}

	return amended
}

func (oa *OrderAmender) publishAmendEvent(ID uuid.UUID, quantity uint) {

	ev := events.NewOrderAmended(ID.String(), quantity, time.Now().UTC(), uint(1))
	env, err := events.NewOrderEventEnvelope(ev, ev.EventType)

	if err != nil {
		log.Printf("Failed to create envelope: %s", err.Error())
	}

	err = oa.publisher.Publish(env)
	if err != nil {
		log.Printf("Failed to publish cancelled event: %s", ev.String())
	}
}

func (oa *OrderAmender) getAmendQuantity(orig uint, amend uint) (uint, error) {
	if orig >= amend {
		return 0, errors.New("Amend quantity less or equal than orders")
	}

	return amend - orig, nil
}

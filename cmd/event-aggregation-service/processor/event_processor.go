package processor

import (
	"github.com/mantzas/incata"
	uuid "github.com/satori/go.uuid"
	"github.com/tradsim/tradsim-go/cmd/event-aggregation-service/aggregator"
	"github.com/tradsim/tradsim-go/cmd/event-aggregation-service/data"
	"github.com/tradsim/tradsim-go/events"
)

// Processor interface
type Processor interface {
	Process(event events.OrderEventStored) error
}

// EventProcessor processes stored events
type EventProcessor struct {
	evr   incata.Retriever
	evagg aggregator.EventAggregator
	oragg aggregator.OrderAggregator
	repo  data.OrderRepository
}

// NewEventProcessor creates a new event processor
func NewEventProcessor(evr incata.Retriever, evagg aggregator.EventAggregator,
	oragg aggregator.OrderAggregator, repo data.OrderRepository) *EventProcessor {
	return &EventProcessor{evr, evagg, oragg, repo}
}

// Process the stored event
func (ep *EventProcessor) Process(event events.OrderEventStored) error {

	sourceID, err := uuid.FromString(event.OrderID)
	if err != nil {
		return err
	}

	events, err := ep.evr.Retrieve(sourceID)
	if err != nil {
		return err
	}

	or, err := ep.evagg.Aggregate(events)
	if err != nil {
		return err
	}

	orders, err := ep.repo.GetOrders()
	if err != nil {
		return err
	}

	_ = ep.oragg.Aggregate(or.Symbol, orders)

	// TODO: Store order and position to db

	return nil
}

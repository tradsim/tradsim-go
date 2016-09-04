package processor

import (
	"github.com/mantzas/incata"
	uuid "github.com/satori/go.uuid"
	"github.com/tradsim/tradsim-go/cmd/event-aggregation-service/aggregator"
	"github.com/tradsim/tradsim-go/events"
)

// Processor interface
type Processor interface {
	Process(event events.OrderEventStored) error
}

// EventProcessor processes stored events
type EventProcessor struct {
	evr incata.Retriever
	agg aggregator.EventAggregator
}

// NewEventProcessor creates a new event processor
func NewEventProcessor(evr incata.Retriever, agg aggregator.EventAggregator) *EventProcessor {
	return &EventProcessor{evr, agg}
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

	_, err = ep.agg.Aggregate(events)
	if err != nil {
		return err
	}

	// TODO: Calculate position

	// TODO: Store order and position to db

	return nil
}

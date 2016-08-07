package events

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetEventTypeError(t *testing.T) {

	require := require.New(t)

	_, err := GetEventType("test")

	require.NotNil(err)
}

var orderEnvelopeGetErrorTests = []struct {
	in OrderEventType
}{
	{OrderCreatedType},
	{OrderAmendedType},
	{OrderCanceledType},
	{OrderTradedType},
}

func TestOrderEventEnvelopeError(t *testing.T) {

	require := require.New(t)

	for _, tt := range orderEnvelopeGetErrorTests {

		envelope := OrderEventEnvelope{tt.in, "TEST"}
		_, err := envelope.GetOrderEvent()
		require.NotNil(err)
	}
}

type input struct {
	data      interface{}
	eventType OrderEventType
}

var orderEnvelopeTests = []struct {
	in  input
	out OrderEventType
}{
	{input{OrderCreated{}, OrderCreatedType}, OrderCreatedType},
	{input{OrderAmended{}, OrderAmendedType}, OrderAmendedType},
	{input{OrderCanceled{}, OrderCanceledType}, OrderCanceledType},
	{input{OrderTraded{}, OrderTradedType}, OrderTradedType},
}

func TestNewOrderEventEnvelope(t *testing.T) {

	require := require.New(t)

	for _, tt := range orderEnvelopeTests {

		envelope, err := NewOrderEventEnvelope(tt.in.data, tt.in.eventType)

		require.Nil(err)
		require.Equal(envelope.EventType, tt.out)

		event, err := envelope.GetOrderEvent()

		require.Nil(err)

		eventType, err := GetEventType(event)

		require.Nil(err)
		require.Equal(tt.out, eventType)
	}
}

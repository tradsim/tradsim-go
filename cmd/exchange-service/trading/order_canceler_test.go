package trading

import (
	"testing"

	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
	"github.com/tradsim/tradsim-go/mocks"
	"github.com/tradsim/tradsim-go/models"
)

func TestNewOrderCanceller(t *testing.T) {

	require := require.New(t)

	canceller := NewOrderCanceller(&mocks.MockPublisher{})

	require.NotNil(canceller)

}

func TestCancelOrder(t *testing.T) {

	require := require.New(t)

	book := models.NewOrderBook()
	order := models.NewOrder(uuid.NewV4(), "TT", 199.99, 10, models.Sell)

	ap := NewOrderAppender()
	ap.Append(book, order)
	publisher := &mocks.MockPublisher{}

	cnc := NewOrderCanceller(publisher)
	cancelled := cnc.Cancel(book, order.ID)

	require.True(cancelled)
	require.Equal(models.Cancelled, order.Status)
	prices, _ := book.Symbols["TT"]
	require.Equal(models.Cancelled, prices[0].Sell.Orders[0].Status)
	require.Len(publisher.Envelopes, 1)
}

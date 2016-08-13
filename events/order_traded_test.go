package events

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestOrderTradeString(t *testing.T) {

	require := require.New(t)

	l, _ := time.LoadLocation("Europe/Athens")
	dt := time.Date(2016, 8, 13, 17, 33, 11, 111, l)
	event := NewOrderTraded("d1de4242-6620-4030-b2a7-4a701631c3ba", dt, 1.99, 10, 1)

	require.Equal("OrderTraded: [d1de4242-6620-4030-b2a7-4a701631c3ba] 2016-08-13 17:33:11.000000111 +0300 EEST 1 10@1.990000", event.String())
}

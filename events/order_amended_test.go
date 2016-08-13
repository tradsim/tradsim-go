package events

import (
	"os"
	"testing"
	"time"

	"github.com/mantzas/adaptlog"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	adaptlog.ConfigureStdLevelLogger(adaptlog.DebugLevel, nil, "")
	retCode := m.Run()
	os.Exit(retCode)
}

func TestOrderAmendedString(t *testing.T) {

	require := require.New(t)
	l, _ := time.LoadLocation("Europe/Athens")
	dt := time.Date(2016, 8, 13, 17, 33, 11, 111, l)

	event := NewOrderAmended("d1de4242-6620-4030-b2a7-4a701631c3ba", 1, dt, 1)

	require.Equal("OrderAmended: [d1de4242-6620-4030-b2a7-4a701631c3ba] 2016-08-13 17:33:11.000000111 +0300 EEST 1 1", event.String())
}

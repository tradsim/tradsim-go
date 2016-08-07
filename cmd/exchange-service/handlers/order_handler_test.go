package handlers

import (
	"testing"

	"os"

	"github.com/mantzas/adaptlog"
)

func TestMain(m *testing.M) {
	adaptlog.ConfigureStdLevelLogger(adaptlog.DebugLevel, nil, "")
	retCode := m.Run()
	os.Exit(retCode)
}

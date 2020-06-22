package strategy_conf_loader

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// setup code...
	code := m.Run()
	// teardown code...
	os.Exit(code)
}

func TestInit(t *testing.T) {
}

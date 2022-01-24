package log

import (
	"testing"
)

func TestLogger(t *testing.T) {
	SetFlags(Lshortfile | Lmsgprefix)
	Info("Infotest")
	Debug("Should not be visible")
	Verbose = true
	Debug("Should be visible")
	Warning("A warning")
}

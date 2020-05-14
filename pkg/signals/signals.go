// Package signals exposes signal handler.
package signals

import (
	"os"
	"os/signal"
)

// Handler defines signal handler.
type Handler struct {
	// SignalCh stores signals received from the system.
	SignalCh chan os.Signal
}

// NewSignalHandler creates a new Handler.
func NewSignalHandler() *Handler {
	// Different signals that we are listening to. For the time being putting
	// only os.Interrupt, making a slice so that we can easily add more.
	signals := []os.Signal{os.Interrupt}

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, signals...)

	return &Handler{
		SignalCh: signalCh,
	}
}

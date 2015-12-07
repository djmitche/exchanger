package exchanger

import (
	"fmt"
)

const (
	InvalidTick = iota

	// Exchanges indicate new limit orders with a quote tick
	QuoteTick

	// Exchanges indicate matching an order with an execution tick
	ExecutionTick

	// Exchanges indicate the withdrawal of a limit order with a cancellation tick
	CancellationTick
)

// Exchanges report changes in their state with Ticks
type Tick struct {
	// The type of the tick (one of the *Tick constants)
	Type int

	// The quantity of contracts concerned
	Quantity int

	// The price of the contracts
	Price int

	// The symbol for the contracts
	Symbol string
}

// A ticker takes incoming Ticks and does whatever it would like with them
type Ticker interface {
	Tick(*Tick)
}

func (tk *Tick) String() string {
	var ty string
	switch tk.Type {
	case QuoteTick:
		ty = "Q"
	case ExecutionTick:
		ty = "E"
	case CancellationTick:
		ty = "C"
	default:
		ty = "???"
	}

	return fmt.Sprintf("<%s %dx%s@%d>", ty, tk.Quantity, tk.Symbol, tk.Price)
}

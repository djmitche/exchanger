package exchanger

import (
	"fmt"
)

const (
	InvalidOrder = iota

	// A limit order is placed on the exchange's book and will be matched against
	// incoming market orders.
	LimitOrder

	// A market order is matched against existing limit orders and either executed
	// or discarded.
	MarketOrder
)

// An Order is a message sent to an order processor (usually an exchange),
// representing a request to change the sender's position in the market
type Order struct {
	// Type of the order (one of the *Order constants)
	Type int

	// The quantity of contracts concerned
	Quantity int

	// The price of the contracts
	Price int

	// The symbol for the contracts
	Symbol string

	// The party-specific ordinal for this order; used, for example, to cancel a limit order
	Ordinal int
}

// An order processor takes incoming orders and does someting appropriate with
// them.  The method must handle all erorrs internally.
type OrderProcessor interface {
	// TODO: result "channel"?
	Process(*Order)
}

func (o *Order) String() string {
	var ty string
	switch o.Type {
	case LimitOrder:
		ty = "LIM"
	case MarketOrder:
		ty = "MKT"
	default:
		ty = "???"
	}
	return fmt.Sprintf("<%s %dx%s@%d>", ty, o.Quantity, o.Symbol, o.Price)
}

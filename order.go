package exchanger

import (
	"fmt"
)

const (
	// A fill event means that the order was (maybe partially) filled
	// at the given price and quantity.
	FillEvent = iota

	// An expire event indicates that the order has expired and been
	// removed from the books.  Price and quantity are zero, and the
	// event is always final.
	ExpireEvent = iota
)

// An OrderEvent is an event relating to a specific order after it has been
// submitted to a processor.
type OrderEvent struct {
	// One of the *Event constants
	EventType int

	// The quantity involved in the transaction
	Quantity int

	// The price of the transaction
	Price int

	// True if this is the last event for this order
	IsFinal bool
}

const (
	Buy   = 1 << iota // 0 = sell
	Limit = 1 << iota // 0 = market
)

// An Order is a message sent to an order processor (usually an exchange),
// representing a request to change the sender's position in the market
type Order struct {
	// Order characteristics (bitfield)
	OrderInfo uint32

	// The quantity of contracts concerned
	Quantity int

	// The price of the contracts
	Price int

	// The symbol for the contracts
	Symbol string

	// Callback for events for this order
	Callback func(OrderEvent)
}

func (o *Order) String() string {
	var buySell, limitMarket string
	if o.IsBuy() {
		buySell = "BUY"
	} else {
		buySell = "SELL"
	}
	if o.IsLimit() {
		limitMarket = "LIM"
	} else {
		limitMarket = "MKT"
	}
	return fmt.Sprintf("<%s/%s %dx%s@%d>", buySell, limitMarket,
		o.Quantity, o.Symbol, o.Price)
}

func (o *Order) IsBuy() bool    { return o.OrderInfo&Buy != 0 }
func (o *Order) IsSell() bool   { return o.OrderInfo&Buy == 0 }
func (o *Order) IsLimit() bool  { return o.OrderInfo&Limit != 0 }
func (o *Order) IsMarket() bool { return o.OrderInfo&Limit == 0 }

// An order processor takes incoming orders and does someting appropriate with
// them.  The method must handle all erorrs internally.
type OrderProcessor interface {
	Process(*Order)
}

package exchange

import (
	"github.com/djmitche/exchanger"
)

type order struct {
	exchanger.Order

	// marker of order arrival; buy orders are negative and sell orders
	// are positive, with ordinals closer to zero being older; this supports
	// the FIFO sort order
	ordinal uint64
}

func (o *order) callback(oe exchanger.OrderEvent) {
	if o.Callback != nil {
		o.Callback(oe)
	}
}

func (o *order) String() string {
	return o.Order.String()
}

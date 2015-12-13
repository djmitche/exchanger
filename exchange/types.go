package exchange

import (
	"github.com/djmitche/exchanger"
)

// TODO: rename to just order
type stampedOrder struct {
	exchanger.Order

	// marker of order arrival; buy orders are negative and sell orders
	// are positive, with ordinals closer to zero being older; this supports
	// the FIFO sort order
	ordinal uint64
}

func (o *stampedOrder) String() string {
	return o.Order.String()
}

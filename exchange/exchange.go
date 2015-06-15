package exchange

import (
	"fmt"
	"strings"
    "sort"
)

type book []Order

type Exchange struct {
	Symbol string

	bids book
	asks book
}

func (e Exchange) String() string {
    fmtOrders := func (orders book) string {
        res := make([]string, len(orders))
        for i, o := range orders {
            res[i] = o.String()
        }
        return fmt.Sprintf("[%s]", strings.Join(res, ", "))

    }

	return fmt.Sprintf("Exchange for %s with bids: %v; asks: %v>",
		e.Symbol, fmtOrders(e.bids), fmtOrders(e.asks))
}

// books are sorted from most to least preferred order; that is by price
// (low to high for asks, high to low for buys) then by ordinal
func (b book) Len() int {
    return len(b)
}

func (b book) Less(i, j int) bool {
    oi := b[i]
    oj := b[j]

    if oi.Price == oj.Price {
        return oi.ordinal < oj.ordinal
    }

    if oi.OrderType == "ASK" {
        return oi.Price < oj.Price
    } else {
        return oi.Price > oj.Price
    }
}

func (b book) Swap(i, j int) {
    b[i], b[j] = b[j], b[i]
}

func (e *Exchange) normalize() {
    // TODO: omg slow
    filter := func (unfiltered book) (filtered book) {
        for _, o := range(unfiltered) {
            if o.Quantity != 0 {
                filtered = append(filtered, o)
            }
        }
        return
    }

    e.bids = filter(e.bids)
    sort.Sort(e.bids)

    e.asks = filter(e.asks)
    sort.Sort(e.asks)
}

func (e *Exchange) handleBid(bid *Order, execs ExecutionChan) {
    defer e.normalize()
    // match against asks, in order, until exhausted
    for i := range(e.asks) {
        exec := match(bid, &e.asks[i])
        if exec != nil {
            execs <- exec
        }

        if bid.Quantity == 0 {
            return
        }
    }
    e.bids = append(e.bids, *bid)
	return
}

func (e *Exchange) handleAsk(ask *Order, execs ExecutionChan) {
    defer e.normalize()
    // match against bids, in order, until exhausted
    for i := range(e.bids) {
        exec := match(&e.bids[i], ask)
        if exec != nil {
            execs <- exec
        }

        if ask.Quantity == 0 {
            return
        }
    }
    e.asks = append(e.asks, *ask)
	return
}

func (e *Exchange) Run(orders OrderChan, execs ExecutionChan) {
    var ordinal int64
    for order := range orders {
        order.ordinal = ordinal
        ordinal += 1

        switch order.OrderType {
        case "BID": e.handleBid(order, execs)
        case "ASK": e.handleAsk(order, execs)
        }
    }
}

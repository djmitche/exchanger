package main

import (
	"fmt"
	"strings"
    "sort"
)

type Book []Order

type Exchange struct {
	symbol string

	bids Book
	asks Book
}

func (e Exchange) String() string {
    fmtOrders := func (orders Book) string {
        res := make([]string, len(orders))
        for i, o := range orders {
            res[i] = o.String()
        }
        return fmt.Sprintf("[%s]", strings.Join(res, ", "))

    }

	return fmt.Sprintf("Exchange for %s with bids: %v; asks: %v>",
		e.symbol, fmtOrders(e.bids), fmtOrders(e.asks))
}

// books are sorted from most to least preferred order; that is by price
// (low to high for asks, high to low for buys) then by ordinal
func (b Book) Len() int {
    return len(b)
}

func (b Book) Less(i, j int) bool {
    oi := b[i]
    oj := b[j]

    if oi.price == oj.price {
        return oi.ordinal < oj.ordinal
    }

    if oi.orderType == "ASK" {
        return oi.price < oj.price
    } else {
        return oi.price > oj.price
    }
}

func (b Book) Swap(i, j int) {
    b[i], b[j] = b[j], b[i]
}

func (e *Exchange) normalize() {
    // TODO: omg slow
    filter := func (unfiltered Book) (filtered Book) {
        for _, o := range(unfiltered) {
            if o.quantity != 0 {
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

        if bid.quantity == 0 {
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

        if ask.quantity == 0 {
            return
        }
    }
    e.asks = append(e.asks, *ask)
	return
}

func (e *Exchange) run(orders OrderChan, execs ExecutionChan) {
    var ordinal int64
    for order := range orders {
        order.ordinal = ordinal
        ordinal += 1

        switch order.orderType {
        case "BID": e.handleBid(order, execs)
        case "ASK": e.handleAsk(order, execs)
        }
    }
}

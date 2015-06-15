package exchange

import (
	"fmt"
	"strings"
    "sort"
)

type book []Order

type market struct {
    symbol string

	bids book
	asks book
}

func (m market) String() string {
    fmtOrders := func (orders book) string {
        res := make([]string, len(orders))
        for i, o := range orders {
            res[i] = o.String()
        }
        return fmt.Sprintf("[%s]", strings.Join(res, ", "))
    }

	return fmt.Sprintf("market in %s with book %v / %v>",
		m.symbol, fmtOrders(m.bids), fmtOrders(m.asks))
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

func (m *market) normalize() {
    // TODO: omg slow
    filter := func (unfiltered book) (filtered book) {
        for _, o := range(unfiltered) {
            if o.Quantity != 0 {
                filtered = append(filtered, o)
            }
        }
        return
    }

    m.bids = filter(m.bids)
    sort.Sort(m.bids)

    m.asks = filter(m.asks)
    sort.Sort(m.asks)
}

func (m *market) handleBid(bid *Order, execs ExecutionChan) {
    defer m.normalize()
    // match against asks, in order, until exhausted
    for i := range(m.asks) {
        exec := match(bid, &m.asks[i])
        if exec != nil {
            execs <- exec
        }

        if bid.Quantity == 0 {
            return
        }
    }
    m.bids = append(m.bids, *bid)
	return
}

func (m *market) handleAsk(ask *Order, execs ExecutionChan) {
    defer m.normalize()
    // match against bids, in order, until exhausted
    for i := range(m.bids) {
        exec := match(&m.bids[i], ask)
        if exec != nil {
            execs <- exec
        }

        if ask.Quantity == 0 {
            return
        }
    }
    m.asks = append(m.asks, *ask)
	return
}

func (m *market) handleOrder(order *Order, execs ExecutionChan) {
    switch order.OrderType {
    case "BID": m.handleBid(order, execs)
    case "ASK": m.handleAsk(order, execs)
    }
}

type Exchange struct {
    markets map[string]*market
}

func (e Exchange) String() string {
    res := make([]string, 0)
    for _, mkt := range e.markets {
        res = append(res, mkt.String())
    }
    return fmt.Sprintf("[%s]", strings.Join(res, ", "))
}

func (e *Exchange) Run(orders OrderChan, execs ExecutionChan) {
    var ordinal int64

    e.markets = make(map[string]*market, 0)

    for order := range orders {
        order.ordinal = ordinal
        ordinal += 1

        mkt, ok := e.markets[order.Symbol]
        if !ok {
            mkt = &market{symbol: order.Symbol}
            e.markets[order.Symbol] = mkt
        }
        mkt.handleOrder(order, execs)
    }
}

package main

import (
	"log"
	"testing"
)

func tryHandle(t *testing.T, bids, asks []Order, order Order, exp_bids, exp_asks []Order, exp_execs []Execution) {
    // save some typing in the test functions and fill in constant details here
    for i := range bids {
        bids[i].symbol = "A"
        bids[i].orderType = "BID"
        bids[i].party = "Bubba"
        bids[i].ordinal = 2 * int64(i)
    }
    for i := range asks {
        asks[i].symbol = "A"
        asks[i].orderType = "ASK"
        asks[i].party = "Sue"
        asks[i].ordinal = 2 * int64(i) + 1
    }
    order.symbol = "A"
    order.party = "Burt"
    order.ordinal = 1000 // after any of the others

    // create and normalize the exchange
	exch := Exchange{symbol: "A", bids: bids, asks: asks}
    exch.normalize()

	log.Printf("before: %s", exch)
    log.Printf("handling %s", order)

    var execs []Execution
    switch order.orderType {
    case "BID": execs = exch.handleBid(&order)
    case "ASK": execs = exch.handleAsk(&order)
    default: t.Fatal("order doesn't have a type")
    }

    for _, exec := range execs {
        log.Printf("executed %s", exec)
    }
	log.Printf("after: %s", exch)

    if len(exch.bids) != len(exp_bids) {
        t.Fatal("unexpected number of bids in book")
    }
    for i, exp := range exp_bids {
        got := exch.bids[i]
        if got.quantity != exp.quantity || got.price != exp.price {
            t.Fatal("unexpected bid in book")
        }
    }

    if len(exch.asks) != len(exp_asks) {
        t.Fatal("unexpected number of asks in book")
    }
    for i, exp := range exp_asks {
        got := exch.asks[i]
        if got.quantity != exp.quantity || got.price != exp.price {
            t.Fatal("unexpected ask in book")
        }
    }

    if len(execs) != len(exp_execs) {
        t.Fatal("unexpected number of executions")
    }
    for i, exp := range exp_execs {
        got := execs[i]
        if got.quantity != exp.quantity || got.price != exp.price {
            t.Fatal("unexpected execution")
        }
    }
}

func TestHandleBid(t *testing.T) {
	tryHandle(t,
		[]Order{},
		[]Order{},
		Order{orderType: "BID", quantity: 10, price: 15},
		[]Order{Order{quantity: 10, price: 15}},
		[]Order{},
        []Execution{},
	)
}

func TestHandleBidExisting(t *testing.T) {
	tryHandle(t,
		[]Order{Order{quantity: 10, price: 13}},
		[]Order{},
		Order{orderType: "BID", quantity: 10, price: 15},
		[]Order{Order{quantity: 10, price: 15}, Order{quantity: 10, price: 13}},
		[]Order{},
        []Execution{},
	)
}

func TestHandleBidNotMatched(t *testing.T) {
	tryHandle(t,
		[]Order{},
		[]Order{Order{quantity: 10, price: 15}},
		Order{orderType: "BID", quantity: 10, price: 13},
		[]Order{Order{quantity: 10, price: 13}},
		[]Order{Order{quantity: 10, price: 15}},
        []Execution{},
	)
}

func TestHandleBidMatched(t *testing.T) {
	tryHandle(t,
		[]Order{},
		[]Order{Order{quantity: 10, price: 13}},
		Order{orderType: "BID", quantity: 10, price: 13},
		[]Order{},
		[]Order{},
        []Execution{Execution{quantity: 10, price: 13}},
	)
}

func TestHandleBidMatchedNotBest(t *testing.T) {
	tryHandle(t,
		[]Order{},
		[]Order{
            Order{quantity: 10, price: 13},
            Order{quantity: 10, price: 14},
            Order{quantity: 10, price: 15},
        },
		Order{orderType: "BID", quantity: 10, price: 14},
		[]Order{},
		[]Order{
            // note the ask at 13 was executed, but at the bid price (14)
            Order{quantity: 10, price: 14},
            Order{quantity: 10, price: 15},
        },
        []Execution{Execution{quantity: 10, price: 14}},
	)
}

func TestHandleBidPartiallyMatchedAsk(t *testing.T) {
	tryHandle(t,
		[]Order{},
		[]Order{Order{quantity: 20, price: 13}},
		Order{orderType: "BID", quantity: 10, price: 13},
		[]Order{},
		[]Order{Order{quantity: 10, price: 13}},
        []Execution{Execution{quantity: 10, price: 13}},
	)
}

func TestHandleBidPartiallyMatchedBid(t *testing.T) {
	tryHandle(t,
		[]Order{},
		[]Order{Order{quantity: 10, price: 13}},
		Order{orderType: "BID", quantity: 20, price: 13},
		[]Order{Order{quantity: 10, price: 13}},
		[]Order{},
        []Execution{Execution{quantity: 10, price: 13}},
	)
}

func TestHandleBidMatchesMultiple(t *testing.T) {
	tryHandle(t,
		[]Order{},
		[]Order{Order{quantity: 5, price: 13}, Order{quantity: 5, price: 13}},
		Order{orderType: "BID", quantity: 15, price: 13},
		[]Order{Order{quantity: 5, price: 13}},
		[]Order{},
        []Execution{Execution{quantity: 5, price: 13}, Execution{quantity: 5, price: 13}},
	)
}

func TestHandleAsk(t *testing.T) {
	tryHandle(t,
		[]Order{},
		[]Order{},
		Order{orderType: "ASK", quantity: 10, price: 15},
		[]Order{},
		[]Order{Order{quantity: 10, price: 15}},
        []Execution{},
	)
}

func TestHandleAskExisting(t *testing.T) {
	tryHandle(t,
		[]Order{},
		[]Order{Order{quantity: 10, price: 13}},
		Order{orderType: "ASK", quantity: 10, price: 15},
		[]Order{},
		[]Order{Order{quantity: 10, price: 13}, Order{quantity: 10, price: 15}},
        []Execution{},
	)
}

func TestHandleAskNotMatched(t *testing.T) {
	tryHandle(t,
		[]Order{Order{quantity: 10, price: 13}},
		[]Order{},
		Order{orderType: "ASK", quantity: 10, price: 15},
		[]Order{Order{quantity: 10, price: 13}},
		[]Order{Order{quantity: 10, price: 15}},
        []Execution{},
	)
}

func TestHandleAskMatched(t *testing.T) {
	tryHandle(t,
		[]Order{Order{quantity: 10, price: 13}},
		[]Order{},
		Order{orderType: "ASK", quantity: 10, price: 13},
		[]Order{},
		[]Order{},
        []Execution{Execution{quantity: 10, price: 13}},
	)
}

func TestHandleAskMatchedNotBest(t *testing.T) {
	tryHandle(t,
		[]Order{
            Order{quantity: 10, price: 15},
            Order{quantity: 10, price: 14},
            Order{quantity: 10, price: 13},
        },
		[]Order{},
		Order{orderType: "ASK", quantity: 10, price: 14},
		[]Order{
            // note the bid at 15 was executed, but at the ask price (14)
            Order{quantity: 10, price: 14},
            Order{quantity: 10, price: 13},
        },
		[]Order{},
        []Execution{Execution{quantity: 10, price: 14}},
	)
}

func TestHandleAskPartiallyMatchedAsk(t *testing.T) {
	tryHandle(t,
		[]Order{Order{quantity: 20, price: 13}},
		[]Order{},
		Order{orderType: "ASK", quantity: 10, price: 13},
		[]Order{Order{quantity: 10, price: 13}},
		[]Order{},
        []Execution{Execution{quantity: 10, price: 13}},
	)
}

func TestHandleAskPartiallyMatchedBid(t *testing.T) {
	tryHandle(t,
		[]Order{Order{quantity: 10, price: 13}},
		[]Order{},
		Order{orderType: "ASK", quantity: 20, price: 13},
		[]Order{},
		[]Order{Order{quantity: 10, price: 13}},
        []Execution{Execution{quantity: 10, price: 13}},
	)
}

func TestHandleAskMatchesMultiple(t *testing.T) {
	tryHandle(t,
		[]Order{Order{quantity: 5, price: 13}, Order{quantity: 5, price: 13}},
		[]Order{},
		Order{orderType: "ASK", quantity: 15, price: 13},
		[]Order{},
		[]Order{Order{quantity: 5, price: 13}},
        []Execution{Execution{quantity: 5, price: 13}, Execution{quantity: 5, price: 13}},
	)
}

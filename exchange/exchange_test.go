package exchange

import (
	"log"
	"testing"
)

func tryHandle(t *testing.T, bids, asks []Order, order Order, exp_bids, exp_asks []Order, exp_execs []Execution) {
    // save some typing in the test functions and fill in constant details here
    for i := range bids {
        bids[i].Symbol = "A"
        bids[i].OrderType = "BID"
        bids[i].Party = "Bubba"
        bids[i].ordinal = 2 * int64(i)
    }
    for i := range asks {
        asks[i].Symbol = "A"
        asks[i].OrderType = "ASK"
        asks[i].Party = "Sue"
        asks[i].ordinal = 2 * int64(i) + 1
    }
    order.Symbol = "A"
    order.Party = "Burt"
    order.ordinal = 1000 // after any of the others

    // create and normalize the market
	mkt := market{symbol: "A", bids: bids, asks: asks}
    mkt.normalize()

	log.Printf("before: %s", mkt)
    log.Printf("handling %s", order)

    execChan := make(ExecutionChan, len(exp_execs) + 10)
    mkt.handleOrder(&order, execChan)
    close(execChan)

    var execs []Execution
    for exec := range execChan {
        log.Printf("executed %s", exec)
        execs = append(execs, *exec)
    }
	log.Printf("after: %s", mkt)

    if len(mkt.bids) != len(exp_bids) {
        t.Fatal("unexpected number of bids in book")
    }
    for i, exp := range exp_bids {
        got := mkt.bids[i]
        if got.Quantity != exp.Quantity || got.Price != exp.Price {
            t.Fatal("unexpected bid in book")
        }
    }

    if len(mkt.asks) != len(exp_asks) {
        t.Fatal("unexpected number of asks in book")
    }
    for i, exp := range exp_asks {
        got := mkt.asks[i]
        if got.Quantity != exp.Quantity || got.Price != exp.Price {
            t.Fatal("unexpected ask in book")
        }
    }

    if len(execs) != len(exp_execs) {
        t.Fatal("unexpected number of executions")
    }
    for i, exp := range exp_execs {
        got := execs[i]
        if got.Quantity != exp.Quantity || got.Price != exp.Price {
            t.Fatal("unexpected execution")
        }
    }
}

func TestHandleBid(t *testing.T) {
	tryHandle(t,
		[]Order{},
		[]Order{},
		Order{OrderType: "BID", Quantity: 10, Price: 15},
		[]Order{Order{Quantity: 10, Price: 15}},
		[]Order{},
        []Execution{},
	)
}

func TestHandleBidExisting(t *testing.T) {
	tryHandle(t,
		[]Order{Order{Quantity: 10, Price: 13}},
		[]Order{},
		Order{OrderType: "BID", Quantity: 10, Price: 15},
		[]Order{Order{Quantity: 10, Price: 15}, Order{Quantity: 10, Price: 13}},
		[]Order{},
        []Execution{},
	)
}

func TestHandleBidNotMatched(t *testing.T) {
	tryHandle(t,
		[]Order{},
		[]Order{Order{Quantity: 10, Price: 15}},
		Order{OrderType: "BID", Quantity: 10, Price: 13},
		[]Order{Order{Quantity: 10, Price: 13}},
		[]Order{Order{Quantity: 10, Price: 15}},
        []Execution{},
	)
}

func TestHandleBidMatched(t *testing.T) {
	tryHandle(t,
		[]Order{},
		[]Order{Order{Quantity: 10, Price: 13}},
		Order{OrderType: "BID", Quantity: 10, Price: 13},
		[]Order{},
		[]Order{},
        []Execution{Execution{Quantity: 10, Price: 13}},
	)
}

func TestHandleBidMatchedNotBest(t *testing.T) {
	tryHandle(t,
		[]Order{},
		[]Order{
            Order{Quantity: 10, Price: 13},
            Order{Quantity: 10, Price: 14},
            Order{Quantity: 10, Price: 15},
        },
		Order{OrderType: "BID", Quantity: 10, Price: 14},
		[]Order{},
		[]Order{
            // note the ask at 13 was executed, but at the bid Price (14)
            Order{Quantity: 10, Price: 14},
            Order{Quantity: 10, Price: 15},
        },
        []Execution{Execution{Quantity: 10, Price: 14}},
	)
}

func TestHandleBidPartiallyMatchedAsk(t *testing.T) {
	tryHandle(t,
		[]Order{},
		[]Order{Order{Quantity: 20, Price: 13}},
		Order{OrderType: "BID", Quantity: 10, Price: 13},
		[]Order{},
		[]Order{Order{Quantity: 10, Price: 13}},
        []Execution{Execution{Quantity: 10, Price: 13}},
	)
}

func TestHandleBidPartiallyMatchedBid(t *testing.T) {
	tryHandle(t,
		[]Order{},
		[]Order{Order{Quantity: 10, Price: 13}},
		Order{OrderType: "BID", Quantity: 20, Price: 13},
		[]Order{Order{Quantity: 10, Price: 13}},
		[]Order{},
        []Execution{Execution{Quantity: 10, Price: 13}},
	)
}

func TestHandleBidMatchesMultiple(t *testing.T) {
	tryHandle(t,
		[]Order{},
		[]Order{Order{Quantity: 5, Price: 13}, Order{Quantity: 5, Price: 13}},
		Order{OrderType: "BID", Quantity: 15, Price: 13},
		[]Order{Order{Quantity: 5, Price: 13}},
		[]Order{},
        []Execution{Execution{Quantity: 5, Price: 13}, Execution{Quantity: 5, Price: 13}},
	)
}

func TestHandleAsk(t *testing.T) {
	tryHandle(t,
		[]Order{},
		[]Order{},
		Order{OrderType: "ASK", Quantity: 10, Price: 15},
		[]Order{},
		[]Order{Order{Quantity: 10, Price: 15}},
        []Execution{},
	)
}

func TestHandleAskExisting(t *testing.T) {
	tryHandle(t,
		[]Order{},
		[]Order{Order{Quantity: 10, Price: 13}},
		Order{OrderType: "ASK", Quantity: 10, Price: 15},
		[]Order{},
		[]Order{Order{Quantity: 10, Price: 13}, Order{Quantity: 10, Price: 15}},
        []Execution{},
	)
}

func TestHandleAskNotMatched(t *testing.T) {
	tryHandle(t,
		[]Order{Order{Quantity: 10, Price: 13}},
		[]Order{},
		Order{OrderType: "ASK", Quantity: 10, Price: 15},
		[]Order{Order{Quantity: 10, Price: 13}},
		[]Order{Order{Quantity: 10, Price: 15}},
        []Execution{},
	)
}

func TestHandleAskMatched(t *testing.T) {
	tryHandle(t,
		[]Order{Order{Quantity: 10, Price: 13}},
		[]Order{},
		Order{OrderType: "ASK", Quantity: 10, Price: 13},
		[]Order{},
		[]Order{},
        []Execution{Execution{Quantity: 10, Price: 13}},
	)
}

func TestHandleAskMatchedNotBest(t *testing.T) {
	tryHandle(t,
		[]Order{
            Order{Quantity: 10, Price: 15},
            Order{Quantity: 10, Price: 14},
            Order{Quantity: 10, Price: 13},
        },
		[]Order{},
		Order{OrderType: "ASK", Quantity: 10, Price: 14},
		[]Order{
            // note the bid at 15 was executed, but at the ask Price (14)
            Order{Quantity: 10, Price: 14},
            Order{Quantity: 10, Price: 13},
        },
		[]Order{},
        []Execution{Execution{Quantity: 10, Price: 14}},
	)
}

func TestHandleAskPartiallyMatchedAsk(t *testing.T) {
	tryHandle(t,
		[]Order{Order{Quantity: 20, Price: 13}},
		[]Order{},
		Order{OrderType: "ASK", Quantity: 10, Price: 13},
		[]Order{Order{Quantity: 10, Price: 13}},
		[]Order{},
        []Execution{Execution{Quantity: 10, Price: 13}},
	)
}

func TestHandleAskPartiallyMatchedBid(t *testing.T) {
	tryHandle(t,
		[]Order{Order{Quantity: 10, Price: 13}},
		[]Order{},
		Order{OrderType: "ASK", Quantity: 20, Price: 13},
		[]Order{},
		[]Order{Order{Quantity: 10, Price: 13}},
        []Execution{Execution{Quantity: 10, Price: 13}},
	)
}

func TestHandleAskMatchesMultiple(t *testing.T) {
	tryHandle(t,
		[]Order{Order{Quantity: 5, Price: 13}, Order{Quantity: 5, Price: 13}},
		[]Order{},
		Order{OrderType: "ASK", Quantity: 15, Price: 13},
		[]Order{},
		[]Order{Order{Quantity: 5, Price: 13}},
        []Execution{Execution{Quantity: 5, Price: 13}, Execution{Quantity: 5, Price: 13}},
	)
}

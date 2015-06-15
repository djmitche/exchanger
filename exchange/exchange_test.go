package exchange

import (
	"log"
	"testing"
)

func assertBook(t *testing.T, got, exp book) {
    if len(got) != len(exp) {
        t.Fatal("unexpected number of orders in book")
    }
    for i, e := range exp {
        g := got[i]
        if g.Quantity != e.Quantity || g.Price != e.Price {
            t.Fatal("unexpected order in book")
        }
    }
}

func TestBidOrdering(t *testing.T) {
    m := market{symbol: "X"}
    execChan := make(ExecutionChan)

    m.handleOrder(&Order{OrderType: "BID", Price: 25, Quantity: 1, ordinal: 1}, execChan)
    m.handleOrder(&Order{OrderType: "BID", Price: 26, Quantity: 1, ordinal: 2}, execChan)
    m.handleOrder(&Order{OrderType: "BID", Price: 27, Quantity: 1, ordinal: 3}, execChan)
    m.handleOrder(&Order{OrderType: "BID", Price: 26, Quantity: 1, ordinal: 4}, execChan)
    m.handleOrder(&Order{OrderType: "BID", Price: 24, Quantity: 1, ordinal: 5}, execChan)

    log.Println(m)
    assertBook(t, m.bids, book{
        Order{OrderType: "BID", Price: 27, Quantity: 1, ordinal: 3},
        Order{OrderType: "BID", Price: 26, Quantity: 1, ordinal: 2},
        Order{OrderType: "BID", Price: 26, Quantity: 1, ordinal: 4},
        Order{OrderType: "BID", Price: 25, Quantity: 1, ordinal: 1},
        Order{OrderType: "BID", Price: 24, Quantity: 1, ordinal: 5},
    })
}

func TestAskOrdering(t *testing.T) {
    m := market{symbol: "X"}
    execChan := make(ExecutionChan)

    m.handleOrder(&Order{OrderType: "ASK", Price: 25, Quantity: 1, ordinal: 1}, execChan)
    m.handleOrder(&Order{OrderType: "ASK", Price: 26, Quantity: 1, ordinal: 2}, execChan)
    m.handleOrder(&Order{OrderType: "ASK", Price: 27, Quantity: 1, ordinal: 3}, execChan)
    m.handleOrder(&Order{OrderType: "ASK", Price: 26, Quantity: 1, ordinal: 4}, execChan)
    m.handleOrder(&Order{OrderType: "ASK", Price: 24, Quantity: 1, ordinal: 5}, execChan)

    log.Println(m)
    assertBook(t, m.asks, book{
        Order{OrderType: "ASK", Price: 24, Quantity: 1, ordinal: 5},
        Order{OrderType: "ASK", Price: 25, Quantity: 1, ordinal: 1},
        Order{OrderType: "ASK", Price: 26, Quantity: 1, ordinal: 2},
        Order{OrderType: "ASK", Price: 26, Quantity: 1, ordinal: 4},
        Order{OrderType: "ASK", Price: 27, Quantity: 1, ordinal: 3},
    })
}

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

    assertBook(t, mkt.bids, exp_bids)
    assertBook(t, mkt.asks, exp_asks)

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

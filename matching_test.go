package main

import (
    "testing"
    "log"
)

func tryMatch(t *testing.T, bid, ask *Order, exp_qty int, exp_price int) {
    orig_bid := *bid
    orig_ask := *ask
    log.Printf("match(%v, %v)", bid, ask)
    ex := match(bid, ask)
    log.Printf(" -> %v, %v, %v", ex, bid, ask)
    if exp_qty == 0 {
        if ex != nil {
            t.Fatal("unexpected execution")
        }
    } else {
        if ex == nil {
            t.Fatal("no execution")
        }
        if ex.buyer != bid.party {
            t.Fatal("exeuction buyer is incorrect")
        }
        if ex.seller != ask.party {
            t.Fatal("exeuction seller is incorrect")
        }
        if ex.quantity != exp_qty {
            t.Fatal("exeuction quantity is incorrect")
        }
        if ex.price != exp_price {
            t.Fatal("exeuction price is incorrect")
        }
        if ex.symbol != bid.symbol {
            t.Fatal("exeuction symbol is incorrect")
        }
        if orig_bid.quantity - ex.quantity != bid.quantity {
            t.Fatal("bid quantity is incorrect")
        }
        if orig_ask.quantity - ex.quantity != ask.quantity {
            t.Fatal("bid quantity is incorrect")
        }
    }
}

func TestExactMatch(t *testing.T) {
    bid := Order{orderType: "BID", symbol: "A", quantity: 100, price: 13, party: "Burt"}
    ask := Order{orderType: "ASK", symbol: "A", quantity: 100, price: 13, party: "Sally"}
    tryMatch(t, &bid, &ask, 100, 13)
}

func TestLeftoverBid(t *testing.T) {
    bid := Order{orderType: "BID", symbol: "A", quantity: 150, price: 13, party: "Burt"}
    ask := Order{orderType: "ASK", symbol: "A", quantity: 100, price: 13, party: "Sally"}
    tryMatch(t, &bid, &ask, 100, 13)
}

func TestAskOverBid(t *testing.T) {
    bid := Order{orderType: "BID", symbol: "A", quantity: 100, price: 13, party: "Burt"}
    ask := Order{orderType: "ASK", symbol: "A", quantity: 100, price: 15, party: "Sally"}
    tryMatch(t, &bid, &ask, 0, 0)
}

func TestBidOverAsk(t *testing.T) {
    ask := Order{orderType: "ASK", symbol: "A", quantity: 100, price: 13, ordinal: 800, party: "Sally"}
    bid := Order{orderType: "BID", symbol: "A", quantity: 100, price: 15, ordinal: 900, party: "Burt"}
    tryMatch(t, &bid, &ask, 100, 15)
}

func TestAskUnderBid(t *testing.T) {
    bid := Order{orderType: "BID", symbol: "A", quantity: 100, price: 15, ordinal: 800, party: "Burt"}
    ask := Order{orderType: "ASK", symbol: "A", quantity: 100, price: 13, ordinal: 900, party: "Sally"}
    tryMatch(t, &bid, &ask, 100, 13)
}

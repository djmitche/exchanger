package exchange

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
        if ex.buyer != bid.Party {
            t.Fatal("exeuction buyer is incorrect")
        }
        if ex.seller != ask.Party {
            t.Fatal("exeuction seller is incorrect")
        }
        if ex.Quantity != exp_qty {
            t.Fatal("exeuction quantity is incorrect")
        }
        if ex.Price != exp_price {
            t.Fatal("exeuction price is incorrect")
        }
        if ex.Symbol != bid.Symbol {
            t.Fatal("exeuction Symbol is incorrect")
        }
        if orig_bid.Quantity - ex.Quantity != bid.Quantity {
            t.Fatal("bid quantity is incorrect")
        }
        if orig_ask.Quantity - ex.Quantity != ask.Quantity {
            t.Fatal("bid quantity is incorrect")
        }
    }
}

func TestExactMatch(t *testing.T) {
    bid := Order{OrderType: "BID", Symbol: "A", Quantity: 100, Price: 13, Party: "Burt"}
    ask := Order{OrderType: "ASK", Symbol: "A", Quantity: 100, Price: 13, Party: "Sally"}
    tryMatch(t, &bid, &ask, 100, 13)
}

func TestLeftoverBid(t *testing.T) {
    bid := Order{OrderType: "BID", Symbol: "A", Quantity: 150, Price: 13, Party: "Burt"}
    ask := Order{OrderType: "ASK", Symbol: "A", Quantity: 100, Price: 13, Party: "Sally"}
    tryMatch(t, &bid, &ask, 100, 13)
}

func TestAskOverBid(t *testing.T) {
    bid := Order{OrderType: "BID", Symbol: "A", Quantity: 100, Price: 13, Party: "Burt"}
    ask := Order{OrderType: "ASK", Symbol: "A", Quantity: 100, Price: 15, Party: "Sally"}
    tryMatch(t, &bid, &ask, 0, 0)
}

func TestBidOverAsk(t *testing.T) {
    ask := Order{OrderType: "ASK", Symbol: "A", Quantity: 100, Price: 13, ordinal: 800, Party: "Sally"}
    bid := Order{OrderType: "BID", Symbol: "A", Quantity: 100, Price: 15, ordinal: 900, Party: "Burt"}
    tryMatch(t, &bid, &ask, 100, 15)
}

func TestAskUnderBid(t *testing.T) {
    bid := Order{OrderType: "BID", Symbol: "A", Quantity: 100, Price: 15, ordinal: 800, Party: "Burt"}
    ask := Order{OrderType: "ASK", Symbol: "A", Quantity: 100, Price: 13, ordinal: 900, Party: "Sally"}
    tryMatch(t, &bid, &ask, 100, 13)
}

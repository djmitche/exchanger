package exchange

import (
	"github.com/djmitche/exchanger"
	"testing"
)

func mktBuy(price int, quantity int, ordinal uint64) *order {
	return &order{
		Order: exchanger.Order{
			OrderInfo: exchanger.Buy,
			Quantity:  quantity,
			Price:     price,
			Symbol:    "AAPL",
		},
		ordinal: ordinal,
	}
}

func mktSell(price int, quantity int, ordinal uint64) *order {
	return &order{
		Order: exchanger.Order{
			OrderInfo: 0,
			Quantity:  quantity,
			Price:     price,
			Symbol:    "AAPL",
		},
		ordinal: ordinal,
	}
}

func limBuy(price int, quantity int, ordinal uint64) *order {
	return &order{
		Order: exchanger.Order{
			OrderInfo: exchanger.Buy | exchanger.Limit,
			Quantity:  quantity,
			Price:     price,
			Symbol:    "AAPL",
		},
		ordinal: ordinal,
	}
}

func limSell(price int, quantity int, ordinal uint64) *order {
	return &order{
		Order: exchanger.Order{
			OrderInfo: exchanger.Limit,
			Quantity:  quantity,
			Price:     price,
			Symbol:    "AAPL",
		},
		ordinal: ordinal,
	}
}

func makeBook(orders ...*order) book {
	return book{orders: orders}
}

func quoteTick(price int, quantity int) exchanger.Tick {
	return exchanger.Tick{
		Type:     exchanger.QuoteTick,
		Price:    price,
		Quantity: quantity,
		Symbol:   "AAPL",
	}
}

func execTick(price int, quantity int) exchanger.Tick {
	return exchanger.Tick{
		Type:     exchanger.ExecutionTick,
		Price:    price,
		Quantity: quantity,
		Symbol:   "AAPL",
	}
}

func makeTicks(ticks ...exchanger.Tick) []exchanger.Tick {
	return ticks
}

func assertEqualBooks(t *testing.T, got *book, exp *book, description string) {
	gotString := got.String()
	expString := exp.String()

	if gotString != expString {
		t.Errorf("%s: books differ; got:\n%s\n---- expected:\n%s\n----", description, gotString, expString)
	}
}

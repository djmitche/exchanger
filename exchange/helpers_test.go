package exchange

import (
	"github.com/djmitche/exchanger"
	"testing"
)

func mktBuy(price int, quantity int, ordinal uint64) *stampedOrder {
	return &stampedOrder{
		Order: exchanger.Order{
			OrderInfo: exchanger.Buy,
			Quantity:  quantity,
			Price:     price,
			Symbol:    "AAPL",
		},
		ordinal: ordinal,
	}
}

func mktSell(price int, quantity int, ordinal uint64) *stampedOrder {
	return &stampedOrder{
		Order: exchanger.Order{
			OrderInfo: 0,
			Quantity:  quantity,
			Price:     price,
			Symbol:    "AAPL",
		},
		ordinal: ordinal,
	}
}

func limBuy(price int, quantity int, ordinal uint64) *stampedOrder {
	return &stampedOrder{
		Order: exchanger.Order{
			OrderInfo: exchanger.Buy | exchanger.Limit,
			Quantity:  quantity,
			Price:     price,
			Symbol:    "AAPL",
		},
		ordinal: ordinal,
	}
}

func limSell(price int, quantity int, ordinal uint64) *stampedOrder {
	return &stampedOrder{
		Order: exchanger.Order{
			OrderInfo: exchanger.Limit,
			Quantity:  quantity,
			Price:     price,
			Symbol:    "AAPL",
		},
		ordinal: ordinal,
	}
}

func makeBook(orders ...*stampedOrder) book {
	return book{orders: orders}
}

func assertEqualBooks(t *testing.T, a book, b book, description string) {
	aString := a.String()
	bString := b.String()

	if aString != bString {
		t.Errorf("%s: books differ:\n%s\n----\n%s\n----", description, aString, bString)
	}
}

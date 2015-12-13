package exchange

import (
	"testing"
)

func TestHandleOrder(t *testing.T) {
	tests := []struct {
		description string
		before      book
		order       *stampedOrder
		after       book
	}{
		{
			"unfulfilled market buy",
			makeBook(
				mktBuy(10, 100, 1),
				mktBuy(20, 100, 2),
				mktBuy(30, 100, 3),
			),
			mktBuy(10, 100, 4),
			makeBook(
				mktBuy(10, 100, 1),
				mktBuy(20, 100, 2),
				mktBuy(30, 100, 3),
			),
		},
	}
	for _, test := range tests {
		mkt := market{"AAPL", test.before}
		mkt.normalize()
		mkt.handleOrder(test.order)
		assertEqualBooks(t, mkt.book, test.after, test.description)
	}
}

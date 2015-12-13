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
			"unfulfilled market buy makes no change",
			makeBook(
				limBuy(10, 100, 1),
			),
			mktBuy(-1, 100, 4),
			makeBook(
				limBuy(10, 100, 1),
			),
		},
		{
			"completely matched market buy removes resting order",
			makeBook(
				limSell(10, 100, 1),
				limSell(11, 100, 2),
			),
			mktBuy(-1, 100, 4),
			makeBook(
				limSell(11, 100, 2),
			),
		},
		{
			"completely matched limit buy at a matching price removes resting order",
			makeBook(
				limSell(10, 100, 1),
				limSell(11, 100, 2),
			),
			limBuy(10, 100, 4),
			makeBook(
				limSell(11, 100, 2),
			),
		},
		{
			"completely matched limit sell at a matching price removes resting order",
			makeBook(
				limBuy(9, 100, 2),
				limBuy(10, 100, 1),
			),
			limSell(10, 100, 4),
			makeBook(
				limBuy(9, 100, 2),
			),
		},
		{
			"completely matched limit buy at a better price removes resting order",
			// TODO: verify that this traded at 10, not 11
			makeBook(
				limSell(10, 100, 1),
			),
			limBuy(11, 100, 4),
			makeBook(),
		},
		{
			"completely matched limit sell at a better price removes resting order",
			// TODO: verify that this traded at 10, not 9
			makeBook(
				limBuy(10, 100, 1),
			),
			limSell(9, 100, 4),
			makeBook(),
		},
		{
			"completely matched market sell removes resting order",
			makeBook(
				limBuy(9, 100, 2),
				limBuy(10, 100, 1),
			),
			mktSell(-1, 100, 4),
			makeBook(
				limBuy(9, 100, 2),
			),
		},
		{
			"partially matched market sell removes resting order, disappears",
			makeBook(
				limBuy(10, 100, 2),
			),
			mktSell(-1, 100, 4),
			makeBook(),
		},
		{
			"completely matched market sell removes part of resting order",
			makeBook(
				limBuy(9, 100, 2),
				limBuy(10, 200, 1),
			),
			mktSell(-1, 100, 4),
			makeBook(
				limBuy(9, 100, 2),
				limBuy(10, 100, 1),
			),
		},
		{
			"unfulfilled limit buy sits in the book",
			makeBook(
				limBuy(10, 100, 1),
			),
			limBuy(20, 100, 4),
			makeBook(
				limBuy(10, 100, 1),
				limBuy(20, 100, 4),
			),
		},
		{
			"big limit buy walks the book then sits in it",
			makeBook(
				limSell(10, 100, 1),
				limSell(11, 100, 1),
				limSell(11, 100, 1),
			),
			limBuy(15, 500, 4),
			makeBook(
				limBuy(15, 200, 4),
			),
		},
	}
	for _, test := range tests {
		t.Log(test.description)
		mkt := market{"AAPL", test.before}
		mkt.normalize()
		mkt.handleOrder(test.order)
		assertEqualBooks(t, mkt.book, test.after, test.description)
	}
}

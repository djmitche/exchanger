package exchange

import (
	"math/rand"
	"testing"
)

func TestSorting(t *testing.T) {
	tests := []struct {
		description string
		book        book
		expOrdinals []uint64
	}{
		{
			"order by price (ascending), even if the market is crossed",
			makeBook(
				limBuy(10, 8, 1),
				limSell(20, 8, 2),
				limBuy(30, 8, 3),
			),
			[]uint64{1, 2, 3},
		},

		{
			"order by type (buy/sell) within a price, and by ordinal for the same type",
			makeBook(
				limBuy(10, 8, 1),
				limBuy(20, 8, 2),
				limSell(20, 8, 3),
				limBuy(20, 8, 4),
				limSell(20, 8, 5),
			),
			[]uint64{1, 4, 2, 3, 5},
		},
	}

	for _, test := range tests {
		// shuffle the input book
		bk := make([]*stampedOrder, len(test.book.orders))
		for i, j := range rand.Perm(len(test.book.orders)) {
			bk[i] = test.book.orders[j]
		}
		test.book.orders = bk

		test.book.Sort()

		ok := true
		for i, _ := range bk {
			if test.book.orders[i].ordinal != test.expOrdinals[i] {
				ok = false
			}
		}
		if !ok {
			var got []uint64
			for _, ord := range test.book.orders {
				got = append(got, ord.ordinal)
			}
			t.Errorf("Unexpected sort order for case #q - got %#v", test.description, got)
		}
	}
}

func TestCrossed(t *testing.T) {
	tests := []struct {
		description string
		book        book
		crossed     bool
	}{
		{
			"empty book",
			makeBook(),
			false,
		},
		{
			"crossed",
			makeBook(
				limBuy(10, 8, 1),
				limSell(20, 8, 2),
				limBuy(30, 8, 3),
			),
			true,
		},
		{
			"not crossed",
			makeBook(
				limBuy(10, 8, 1),
				limSell(30, 8, 2),
				limBuy(20, 8, 3),
			),
			false,
		},
		{
			"all buys",
			makeBook(
				limBuy(10, 8, 1),
				limBuy(30, 8, 3),
			),
			false,
		},
		{
			"all sells",
			makeBook(
				limSell(10, 8, 1),
				limSell(30, 8, 3),
			),
			false,
		},
	}

	for _, test := range tests {
		t.Log(test.description)
		test.book.Sort()
		crossed := test.book.IsCrossed()
		if (crossed && !test.crossed) || (!crossed && test.crossed) {
			t.Errorf("%q failed", test.description)
		}
	}
}

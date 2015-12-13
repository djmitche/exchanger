package exchange

import (
	"fmt"
	"github.com/djmitche/exchanger"
	"sort"
	"strings"
)

type book struct {
	orders []*stampedOrder
}

func (bk book) Len() int {
	return len(bk.orders)
}

func (bk book) Less(i, j int) bool {
	oi := bk.orders[i]
	oj := bk.orders[j]

	// if prices are the same, then buys are sorted by ordinal; this puts older
	// buy orders later in the list, but older sell orders earlier, conforming
	// to the expected FIFO behavior.
	if oi.Price != oj.Price {
		// cheaper is less than dearer
		return oi.Price < oj.Price
	}

	if oi.OrderInfo&exchanger.Buy != oj.OrderInfo&exchanger.Buy {
		// buy (1) is less than sell (0)
		return oi.OrderInfo&exchanger.Buy > oj.OrderInfo&exchanger.Buy
	}

	if oi.IsBuy() {
		// for buys, later is less than earlier
		return oi.ordinal > oj.ordinal
	} else {
		// for sells, earlier is less than later
		return oi.ordinal < oj.ordinal
	}
}

func (bk book) Swap(i, j int) {
	bk.orders[i], bk.orders[j] = bk.orders[j], bk.orders[i]
}

// Sort the orders in the book; this orders first by ascending price, then buy
// and sell, then by ordinal such orders with smaller ordinals are closer to
// the midprice.
func (bk book) Sort() {
	sort.Sort(bk)
}

func (bk book) String() string {
	// TODO: use tabular output
	var lines []string
	lines = append(lines, fmt.Sprintf("#%4s %1s %4s %5s %4s",
		"ord", "t", "qty", "price", "qty"))
	for _, ord := range bk.orders {
		if ord.IsBuy() {
			lines = append(lines, fmt.Sprintf("#%4d %1s %4d %5d",
				ord.ordinal, "B", ord.Quantity, ord.Price))
		} else {
			lines = append(lines, fmt.Sprintf("#%4d %1s %4s %5d %4d",
				ord.ordinal, "S", "", ord.Price, ord.Quantity))
		}
	}
	return strings.Join(lines, "\n")
}

// Given a sorted book, return true if the book is "crossed", meaning that it
// contains a sell order at a lower price than a buy order.
func (bk book) IsCrossed() bool {
	var seenSell bool

	for _, ord := range bk.orders {
		if seenSell && ord.IsBuy() {
			return true
		}
		if ord.IsSell() {
			seenSell = true
		}
	}

	return false
}

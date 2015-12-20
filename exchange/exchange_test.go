package exchange

import (
	"github.com/djmitche/exchanger"
	"testing"
)

type recordingTicker struct {
	ticks []exchanger.Tick
	t     *testing.T
}

func (tkr *recordingTicker) Tick(tick *exchanger.Tick) {
	tkr.ticks = append(tkr.ticks, *tick)
}

func (tkr *recordingTicker) assertTicks(ticks []exchanger.Tick) {
	different := false
	if len(tkr.ticks) != len(ticks) {
		different = true
	} else {
		for i := range tkr.ticks {
			if tkr.ticks[i] != ticks[i] {
				different = true
			}
		}
	}

	if different {
		tkr.t.Errorf("expected ticks %#v; got ticks %#v", ticks, tkr.ticks)
	}
}

func TestHandleOrder(t *testing.T) {
	tests := []struct {
		description string
		before      book
		order       *stampedOrder
		after       book
		ticks       []exchanger.Tick
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
			makeTicks(),
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
			makeTicks(
				execTick(10, 100),
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
			makeTicks(
				execTick(10, 100),
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
			makeTicks(
				execTick(10, 100),
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
			makeTicks(
				execTick(10, 100),
			),
		},
		{
			"completely matched limit sell at a better price removes resting order",
			// TODO: verify that this traded at 10, not 9
			makeBook(
				limBuy(10, 100, 1),
			),
			limSell(9, 100, 4),
			makeBook(),
			makeTicks(
				execTick(10, 100),
			),
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
			makeTicks(
				execTick(10, 100),
			),
		},
		{
			"partially matched market sell removes resting order, disappears",
			makeBook(
				limBuy(10, 100, 2),
			),
			mktSell(-1, 100, 4),
			makeBook(),
			makeTicks(
				execTick(10, 100),
			),
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
			makeTicks(
				execTick(10, 100),
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
				limBuy(20, 100, 10),
			),
			makeTicks(
				quoteTick(20, 100),
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
				limBuy(15, 200, 10),
			),
			makeTicks(
				execTick(10, 100),
				execTick(11, 100),
				execTick(11, 100),
				quoteTick(15, 200),
			),
		},
	}

	for _, test := range tests {
		t.Log(test.description)
		exch := New([]string{"AAPL"})
		exch.ordinal = 10
		before := test.before // make a copy
		exch.books["AAPL"] = &before
		ticker := recordingTicker{t: t}
		exch.ticker = &ticker
		exch.normalize()
		exch.Process(&test.order.Order)
		assertEqualBooks(t, exch.books["AAPL"], &test.after, test.description)
		ticker.assertTicks(test.ticks)
	}
}

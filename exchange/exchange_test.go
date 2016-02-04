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
		order       *order
		after       book
		ticks       []exchanger.Tick
		events      []event
	}{
		{
			"unfulfilled market buy makes no change, expires",
			makeBook(
				limBuy(10, 100, 1),
			),
			mktBuy(-1, 100, 4),
			makeBook(
				limBuy(10, 100, 1),
			),
			makeTicks(),
			makeEvents(
				makeEvent(4, exchanger.ExpireEvent, 0, 0, true),
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
			makeTicks(
				execTick(10, 100),
			),
			makeEvents(
				makeEvent(4, exchanger.FillEvent, 100, 10, true),
				makeEvent(1, exchanger.FillEvent, 100, 10, true),
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
			makeEvents(
				makeEvent(4, exchanger.FillEvent, 100, 10, true),
				makeEvent(1, exchanger.FillEvent, 100, 10, true),
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
			makeEvents(
				makeEvent(4, exchanger.FillEvent, 100, 10, true),
				makeEvent(1, exchanger.FillEvent, 100, 10, true),
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
			makeEvents(
				makeEvent(4, exchanger.FillEvent, 100, 10, true),
				makeEvent(1, exchanger.FillEvent, 100, 10, true),
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
			makeEvents(
				makeEvent(4, exchanger.FillEvent, 100, 10, true),
				makeEvent(1, exchanger.FillEvent, 100, 10, true),
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
			makeEvents(
				makeEvent(4, exchanger.FillEvent, 100, 10, true),
				makeEvent(1, exchanger.FillEvent, 100, 10, true),
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
			makeEvents(
				makeEvent(4, exchanger.FillEvent, 100, 10, true),
				makeEvent(2, exchanger.FillEvent, 100, 10, true),
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
			makeEvents(
				makeEvent(4, exchanger.FillEvent, 100, 10, true),
				makeEvent(1, exchanger.FillEvent, 100, 10, false),
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
			makeEvents(),
		},
		{
			"big limit buy walks the book then sits in it",
			makeBook(
				limSell(10, 100, 1),
				limSell(11, 100, 2),
				limSell(11, 100, 3),
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
			makeEvents(
				makeEvent(4, exchanger.FillEvent, 100, 10, false),
				makeEvent(1, exchanger.FillEvent, 100, 10, true),
				makeEvent(4, exchanger.FillEvent, 100, 11, false),
				makeEvent(2, exchanger.FillEvent, 100, 11, true),
				makeEvent(4, exchanger.FillEvent, 100, 11, false),
				makeEvent(3, exchanger.FillEvent, 100, 11, true),
			),
		},
	}

	for _, test := range tests {
		var gotEvents []event

		t.Log(test.description)

		// set up a new exchange with the "before" book
		exch := New([]string{"AAPL"})
		exch.ordinal = 10
		before := test.before // make a copy
		exch.books["AAPL"] = &before

		// attach event callbacks to all resting orders
		for i := range before.orders {
			var order = before.orders[i]
			order.Callback = func(oe exchanger.OrderEvent) {
				gotEvents = append(gotEvents, event{
					ordinal:    order.ordinal,
					OrderEvent: oe,
				})
			}
		}

		test.order.Callback = func(oe exchanger.OrderEvent) {
			gotEvents = append(gotEvents, event{
				ordinal:    test.order.ordinal,
				OrderEvent: oe,
			})
		}

		// set up a ticker
		ticker := recordingTicker{t: t}
		exch.ticker = &ticker

		// normalize the book and process the order
		exch.normalize()
		exch.Process(&test.order.Order)

		// verify results
		assertEqualBooks(t, exch.books["AAPL"], &test.after, test.description)
		ticker.assertTicks(test.ticks)
		assertEqualEvents(t, gotEvents, test.events)
	}
}

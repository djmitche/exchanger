package exchange

import (
	"fmt"
	"github.com/djmitche/exchanger"
	"strings"
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

type event struct {
	ordinal uint64
	exchanger.OrderEvent
}

func (e event) String() string {
	var eventType string
	switch e.EventType {
	case exchanger.FillEvent:
		eventType = "FillEvent"
	case exchanger.ExpireEvent:
		eventType = "ExpireEvent"
	}

	final := "final"
	if !e.IsFinal {
		final = "not final"
	}
	return fmt.Sprintf("[#%v: %s %d@%d %s]",
		e.ordinal, eventType, e.Quantity, e.Price, final)
}

type eventlist []event

func (el eventlist) String() string {
	var rv []string
	for i := range el {
		rv = append(rv, el[i].String())
	}
	return strings.Join(rv, "\n")
}

func makeEvent(ordinal uint64, eventType int, quantity int, price int, isFinal bool) event {
	return event{
		ordinal: ordinal,
		OrderEvent: exchanger.OrderEvent{
			EventType: eventType,
			Quantity:  quantity,
			Price:     price,
			IsFinal:   isFinal,
		},
	}
}

func makeEvents(events ...event) []event {
	return eventlist(events)
}

func assertEqualBooks(t *testing.T, got *book, exp *book, description string) {
	gotString := got.String()
	expString := exp.String()

	if gotString != expString {
		t.Errorf("%s: books differ; got:\n%s\n---- expected:\n%s\n----", description, gotString, expString)
	}
}

func assertEqualEvents(t *testing.T, got eventlist, exp eventlist) {
	different := false
	if len(got) != len(exp) {
		different = true
	} else {
		for i := range got {
			if got[i] != exp[i] {
				different = true
				break
			}
		}
	}
	if different {
		t.Errorf("events differ; got:\n%s\n---- expected:\n%s\n----", got, exp)
	}
}

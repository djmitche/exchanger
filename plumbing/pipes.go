package plumbing

import (
	"github.com/djmitche/exchanger"
)

// An TickTee "splits" a stream of ticks, sending each tick to all attached
// Tickers
type TickTee struct {
	Tickers []exchanger.Ticker
}

func (tt *TickTee) Tick(tick *exchanger.Tick) {
	for _, t := range tt.Tickers {
		t.Tick(tick)
	}
}

// An OrderTee "splits" a stream of orders, sending each order to all attached
// OrderProcessors
type OrderTee struct {
	Processors []exchanger.OrderProcessor
}

func (ot *OrderTee) Process(order *exchanger.Order) {
	for _, p := range ot.Processors {
		p.Process(order)
	}
}

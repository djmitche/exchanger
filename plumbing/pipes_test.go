package plumbing_test

import (
	"github.com/djmitche/exchanger"
	"github.com/djmitche/exchanger/plumbing"
	"log"
	"os"
)

func ExampleTickTee() {
	logger := log.New(os.Stdout, "", 0)
	tee := plumbing.TickTee{
		Tickers: []exchanger.Ticker{
			&plumbing.TickLogger{Logger: logger},
			&plumbing.TickLogger{Logger: logger},
		},
	}
	tee.Tick(&exchanger.Tick{
		Type:     exchanger.QuoteTick,
		Quantity: 100,
		Price:    17,
		Symbol:   "EXCH",
	})
	// Output:
	// <Q 100xEXCH@17>
	// <Q 100xEXCH@17>
}

func ExampleOrderTee() {
	logger := log.New(os.Stdout, "", 0)
	tee := plumbing.OrderTee{
		Processors: []exchanger.OrderProcessor{
			&plumbing.OrderLogger{Logger: logger},
			&plumbing.OrderLogger{Logger: logger},
		},
	}
	tee.Process(&exchanger.Order{
		OrderInfo: 0,
		Quantity:  1000,
		Price:     22,
		Symbol:    "TEST",
	})
	// Output:
	// <SELL/MKT 1000xTEST@22>
	// <SELL/MKT 1000xTEST@22>
}

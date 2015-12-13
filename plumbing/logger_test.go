package plumbing_test

import (
	"github.com/djmitche/exchanger"
	"github.com/djmitche/exchanger/plumbing"
	"log"
	"os"
)

func ExampleTickLogger() {
	logger := plumbing.TickLogger{Logger: log.New(os.Stdout, "", 0)}
	logger.Tick(&exchanger.Tick{
		Type:     exchanger.QuoteTick,
		Quantity: 100,
		Price:    17,
		Symbol:   "EXCH",
	})
	// Output:
	// <Q 100xEXCH@17>
}

func ExampleOrderLogger() {
	logger := plumbing.OrderLogger{Logger: log.New(os.Stdout, "", 0)}
	logger.Process(&exchanger.Order{
		OrderInfo: exchanger.Limit | exchanger.Buy,
		Quantity:  1000,
		Price:     22,
		Symbol:    "TEST",
	})
	// Output:
	// <BUY/LIM 1000xTEST@22>
}

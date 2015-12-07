package plumbing

import (
	"github.com/djmitche/exchanger"
	"log"
)

// A TickLogger is a Ticker that will log every tick it receives
type TickLogger struct {
	Logger *log.Logger
}

func (l *TickLogger) Tick(tick *exchanger.Tick) {
	var println = log.Println
	if l.Logger != nil {
		println = l.Logger.Println
	}
	println(tick)
}

// An OrderLogger is an OrderProcessor that will log every order it receives
type OrderLogger struct {
	Logger *log.Logger
}

func (l *OrderLogger) Process(order *exchanger.Order) {
	var println = log.Println
	if l.Logger != nil {
		println = l.Logger.Println
	}
	println(order)
}

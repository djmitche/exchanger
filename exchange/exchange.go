package exchange

import (
	"fmt"
	"github.com/djmitche/exchanger"
	"log"
	"strings"
)

// An exchange holds a book for each contract and matches orders as they come in
type Exchange struct {
	markets map[string]*market // TODO: use an interface here to allow DI for testing
	log     *log.Logger
	ordinal uint64
}

func New(symbols []string) *Exchange {
	e := Exchange{
		markets: make(map[string]*market),
	}
	for _, sym := range symbols {
		e.markets[sym] = &market{}
	}
	return &e
}

func (e Exchange) String() string {
	res := make([]string, 0)
	for _, mkt := range e.markets {
		res = append(res, mkt.String())
	}
	return fmt.Sprintf("[%s]", strings.Join(res, ", "))
}

func (e *Exchange) Process(o *exchanger.Order) {
	order := &stampedOrder{
		Order:   *o,
		ordinal: e.ordinal,
	}
	e.ordinal++

	mkt, ok := e.markets[order.Symbol]
	if !ok {
		e.log.Printf("Order for nonexistent symbol %s", order.Symbol)
		return
	}
	mkt.handleOrder(order)
}

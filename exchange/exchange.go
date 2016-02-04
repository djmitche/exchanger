package exchange

import (
	"fmt"
	"github.com/djmitche/exchanger"
	"log"
	"strings"
)

// An exchange holds a book for each contract and matches orders as they come in
type Exchange struct {
	books   map[string]*book
	log     *log.Logger
	ordinal uint64
	ticker  exchanger.Ticker
}

func (e *Exchange) normalize() {
	for sym, bk := range e.books {
		// TODO: omg slow!
		var newbook = book{}

		// filter out empty orders
		for _, o := range bk.orders {
			if o.Quantity != 0 {
				newbook.orders = append(newbook.orders, o)
			}
		}

		// sort appropriately
		newbook.Sort()

		// a single book should never be allowed to cross (a limit order that
		// would cross should execute immediately instead)
		if newbook.IsCrossed() {
			panic("book crossed")
		}

		e.books[sym] = &newbook
	}
}

// Match returns true and the matched price if the given aggressive and resting
// order match.  This assumes the caller has verified the two are not both of
// the same type (buy/buy or sell/sell) and of course the same symbol.
func match(aggressive, resting *order) (bool, int) {
	if aggressive.IsMarket() {
		// a market order matches any resting order regardless of price
		return true, resting.Price
	} else {
		if aggressive.IsBuy() {
			if aggressive.Price >= resting.Price {
				return true, resting.Price
			}
		} else {
			if aggressive.Price <= resting.Price {
				return true, resting.Price
			}
		}
	}

	return false, -1
}

// execute the given aggressive and resting orders against one another, assuming they
// match.  The executed quantity is the minimum of the orders' quantities.
func execute(aggressive, resting *order, price int, ticker exchanger.Ticker) {
	quantity := aggressive.Quantity
	if resting.Quantity < quantity {
		quantity = resting.Quantity
	}

	if quantity == 0 {
		return
	}

	// send events and execute the order
	ticker.Tick(&exchanger.Tick{
		Type:     exchanger.ExecutionTick,
		Price:    price,
		Quantity: quantity,
		Symbol:   resting.Symbol,
	})

	event := exchanger.OrderEvent{
		EventType: exchanger.FillEvent,
		Quantity:  quantity,
		Price:     price,
	}

	event.IsFinal = aggressive.Quantity == quantity
	aggressive.callback(event)
	aggressive.Quantity -= quantity

	event.IsFinal = resting.Quantity == quantity
	resting.callback(event)
	resting.Quantity -= quantity
}

func New(symbols []string) *Exchange {
	e := Exchange{
		books: make(map[string]*book),
	}
	for _, sym := range symbols {
		e.books[sym] = &book{}
	}
	return &e
}

func (e Exchange) String() string {
	res := make([]string, 0)
	for _, book := range e.books {
		res = append(res, book.String())
	}
	return fmt.Sprintf("[%s]", strings.Join(res, ", "))
}

func (e *Exchange) Process(o *exchanger.Order) {
	order := &order{
		Order:   *o,
		ordinal: e.ordinal,
	}
	e.ordinal++

	book, ok := e.books[order.Symbol]
	if !ok {
		e.log.Printf("Order for nonexistent symbol %s", order.Symbol)
		return
	}

	var i int
	bookSize := book.Len()

	// seek past the last buy order in the book
	// TODO: make this a book method
	for i = 0; i < bookSize; i++ {
		if book.orders[i].IsSell() {
			break
		}
	}

	if order.IsBuy() {
		// for a buy order, we want to seek upward through the sell side of the book
		for i < bookSize && order.Quantity >= 0 {
			resting := book.orders[i]
			if matched, price := match(order, resting); matched {
				execute(order, resting, price, e.ticker)
			}
			if resting.Quantity != 0 {
				break
			}
			i++
		}
	} else {
		// for a sell order, we want to seek downward through the buy side of the book
		i--
		for i >= 0 && order.Quantity >= 0 {
			resting := book.orders[i]
			if matched, price := match(order, resting); matched {
				execute(order, resting, price, e.ticker)
			}
			if resting.Quantity != 0 {
				break
			}
			i--
		}
	}

	// if this isn't a market order and it isn't filled, add it to the book
	if order.Quantity != 0 {
		if order.IsMarket() {
			order.callback(exchanger.OrderEvent{
				EventType: exchanger.ExpireEvent,
				IsFinal:   true,
			})
		} else {
			fmt.Printf("%#v\n", order)
			e.ticker.Tick(&exchanger.Tick{
				Type:     exchanger.QuoteTick,
				Price:    order.Price,
				Quantity: order.Quantity,
				Symbol:   order.Symbol,
			})
			book.Add(order)
			book, ok = e.books[order.Symbol]
		}
	}

	e.normalize()
}

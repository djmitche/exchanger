package exchange

import (
	"fmt"
	"github.com/djmitche/exchanger"
)

type market struct {
	symbol string
	book   book
}

func (m market) String() string {
	return fmt.Sprintf("%s:\n%s", m.symbol, m.book)
}

func (m *market) normalize() {
	// TODO: omg slow!
	var filtered = book{}
	for _, o := range m.book.orders {
		if o.Quantity != 0 {
			filtered.orders = append(filtered.orders, o)
		}
	}
	m.book = filtered
	m.book.Sort()

	// a single book should never be allowed to cross
	if m.book.IsCrossed() {
		panic("book crossed")
	}
}

// Match returns true and the matched price if the given aggressive and resting
// order match.  This assumes the caller has verified the two are not both of
// the same type (buy/buy or sell/sell).

func match(aggressive, resting *stampedOrder) (bool, int) {
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

func execute(aggressive, resting *stampedOrder, price int, m *market, ticker exchanger.Ticker) {
	quantity := aggressive.Quantity
	if resting.Quantity < quantity {
		quantity = resting.Quantity
	}

	if quantity == 0 {
		return
	}

	// actually execute the order by decrementing quantity
	ticker.Tick(&exchanger.Tick{
		Type:     exchanger.ExecutionTick,
		Price:    price,
		Quantity: quantity,
		Symbol:   m.symbol,
	})
	aggressive.Quantity -= quantity
	resting.Quantity -= quantity
}

func (m *market) handleOrder(order *stampedOrder, ticker exchanger.Ticker) {
	var i int
	bookSize := m.book.Len()

	// seek past the last buy order in the book
	for i = 0; i < bookSize; i++ {
		if m.book.orders[i].IsSell() {
			break
		}
	}

	if order.IsBuy() {
		// for a buy order, we want to seek upward through the sell side of the
		// book
		for i < bookSize && order.Quantity >= 0 {
			resting := m.book.orders[i]
			if matched, price := match(order, resting); matched {
				execute(order, resting, price, m, ticker)
			}
			if resting.Quantity != 0 {
				break
			}
			i++
		}
	} else {
		// for a sell order, we want to seek downward through the buy side of the
		// book
		i--
		for i >= 0 && order.Quantity >= 0 {
			resting := m.book.orders[i]
			if matched, price := match(order, resting); matched {
				execute(order, resting, price, m, ticker)
			}
			if resting.Quantity != 0 {
				break
			}
			i--
		}
	}

	// if this isn't a market order and it isn't filled, add it to the book
	if !order.IsMarket() && order.Quantity != 0 {
		ticker.Tick(&exchanger.Tick{
			Type:     exchanger.QuoteTick,
			Price:    order.Price,
			Quantity: order.Quantity,
			Symbol:   order.Symbol,
		})
		m.book.Add(order)
	}

	m.normalize()
}

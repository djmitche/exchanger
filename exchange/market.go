package exchange

import (
	"fmt"
	//"github.com/djmitche/exchanger"
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
	var filtered = m.book
	for _, o := range m.book.orders {
		if o.Quantity != 0 {
			filtered.orders = append(filtered.orders, o)
		}
	}
	m.book.Sort()

	// a single book should never be allowed to cross
	if m.book.IsCrossed() {
		panic("book crossed")
	}
}

func (m *market) handleOrder(order *stampedOrder) {
}

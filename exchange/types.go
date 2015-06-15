package exchange

import (
	"fmt"
)

// Orders

type Order struct {
	Party     string
	OrderType string
	Quantity  int
	Price     int
	Symbol    string
	ordinal   int64
}

func (o Order) String() string {
    return fmt.Sprintf("<%s \"%s\" %dx%s@%d #%d>", o.OrderType, o.Party, o.Quantity,
        o.Symbol, o.Price, o.ordinal)
}

type OrderChan chan *Order

// Executions

type Execution struct {
    buyer string
    seller string
    Quantity int
    Price int
    Symbol string
}

func (ex Execution) String() string {
    return fmt.Sprintf("<EXEC %dx%s@%d %s to %s>", ex.Quantity, ex.Symbol, ex.Price, ex.seller, ex.buyer)
}

type ExecutionChan chan *Execution

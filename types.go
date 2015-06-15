package main

import (
	"fmt"
)

// Orders

type Order struct {
	party     string
	orderType string
	quantity  int
	price     int
	symbol    string
	ordinal   int64
}

func (o Order) String() string {
    return fmt.Sprintf("<%s \"%s\" %dx%s@%d #%d>", o.orderType, o.party, o.quantity,
        o.symbol, o.price, o.ordinal)
}

type OrderChan chan *Order

// Executions

type Execution struct {
    buyer string
    seller string
    quantity int
    price int
    symbol string
}

func (ex Execution) String() string {
    return fmt.Sprintf("<EXEC %dx%s@%d %s to %s>", ex.quantity, ex.symbol, ex.price, ex.seller, ex.buyer)
}

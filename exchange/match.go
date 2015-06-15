package exchange

// match the given bid and order
func match(bid, ask *Order) (ex *Execution) {
	var price int
    var quantity int

	if bid.Price < ask.Price {
		return // no deal!
    }

    // earlier order is the market-maker, so the later order sets the price, even
    // if that puts the later order at a disadvantage
    if bid.ordinal > ask.ordinal {
        price = bid.Price
    } else {
        price = ask.Price
    }

    // quantity is the minimum of the two orders
	if bid.Quantity <= ask.Quantity {
        quantity = bid.Quantity
    } else {
        quantity = ask.Quantity
    }

    // create an execution and modify the two orders
    ex = &Execution{
        buyer: bid.Party,
        seller: ask.Party,
        Quantity: quantity,
        Price: price,
        Symbol: bid.Symbol,
    }
    bid.Quantity -= quantity
    ask.Quantity -= quantity

    return
}

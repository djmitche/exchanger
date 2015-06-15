package main

// match the given bid and order
func match(bid, ask *Order) (ex *Execution) {
	var price int
    var quantity int

	if bid.price < ask.price {
		return // no deal!
    }

    // earlier order is the market-maker, so the later order sets the price, even
    // if that puts the later order at a disadvantage
    if bid.ordinal > ask.ordinal {
        price = bid.price
    } else {
        price = ask.price
    }

    // quantity is the minimum of the two orders
	if bid.quantity <= ask.quantity {
        quantity = bid.quantity
    } else {
        quantity = ask.quantity
    }

    // create an execution and modify the two orders
    ex = &Execution{
        buyer: bid.party,
        seller: ask.party,
        quantity: quantity,
        price: price,
        symbol: bid.symbol,
    }
    bid.quantity -= quantity
    ask.quantity -= quantity

    return
}

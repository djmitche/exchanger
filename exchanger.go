package main

import (
    "fmt"
    "time"
)

func execPrinter(execs ExecutionChan) {
    for exec := range execs {
        fmt.Println(exec)
    }
}

func main() {
    exch := Exchange{symbol: "A"}
    execs := make(ExecutionChan)
    orders := make(OrderChan)

    go execPrinter(execs)
    go exch.run(orders, execs)

    orders <- &Order{orderType: "BID", party: "Bruce", quantity: 100, price: 92, symbol: "A"}
    orders <- &Order{orderType: "ASK", party: "Sam", quantity: 90, price: 97, symbol: "A"}
    orders <- &Order{orderType: "BID", party: "Bob", quantity: 100, price: 93, symbol: "A"}
    orders <- &Order{orderType: "ASK", party: "Sarah", quantity: 20, price: 94, symbol: "A"}
    orders <- &Order{orderType: "BID", party: "Brian", quantity: 100, price: 91, symbol: "A"}
    orders <- &Order{orderType: "ASK", party: "Samantha", quantity: 20, price: 95, symbol: "A"}
    fmt.Println(exch)
    orders <- &Order{orderType: "BID", party: "Bart", quantity: 100, price: 94, symbol: "A"}
    orders <- &Order{orderType: "BID", party: "Bart", quantity: 100, price: 96, symbol: "A"}

    time.Sleep(1)
}

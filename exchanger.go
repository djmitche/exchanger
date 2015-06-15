package main

import (
    "fmt"
    "time"
    "github.com/djmitche/exchanger/exchange"
)

func execPrinter(execs exchange.ExecutionChan) {
    for exec := range execs {
        fmt.Println(exec)
    }
}

func main() {
    exch := exchange.Exchange{Symbol: "A"}
    execs := make(exchange.ExecutionChan)
    orders := make(exchange.OrderChan)

    go execPrinter(execs)
    go exch.Run(orders, execs)

    orders <- &exchange.Order{OrderType: "BID", Party: "Bruce", Quantity: 100, Price: 92, Symbol: "A"}
    orders <- &exchange.Order{OrderType: "ASK", Party: "Sam", Quantity: 90, Price: 97, Symbol: "A"}
    orders <- &exchange.Order{OrderType: "BID", Party: "Bob", Quantity: 100, Price: 93, Symbol: "A"}
    orders <- &exchange.Order{OrderType: "ASK", Party: "Sarah", Quantity: 20, Price: 94, Symbol: "A"}
    orders <- &exchange.Order{OrderType: "BID", Party: "Brian", Quantity: 100, Price: 91, Symbol: "A"}
    orders <- &exchange.Order{OrderType: "ASK", Party: "Samantha", Quantity: 20, Price: 95, Symbol: "A"}
    fmt.Println(exch)
    orders <- &exchange.Order{OrderType: "BID", Party: "Bart", Quantity: 100, Price: 94, Symbol: "A"}
    orders <- &exchange.Order{OrderType: "BID", Party: "Bart", Quantity: 100, Price: 96, Symbol: "A"}

    time.Sleep(1)
}

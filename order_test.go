package exchanger_test

import (
	"github.com/djmitche/exchanger"
	"testing"
)

func TestOrderString(t *testing.T) {
	tests := []struct {
		input exchanger.Order
		exp   string
	}{
		{exchanger.Order{}, "<??? 0x@0>"},
		{exchanger.Order{Type: exchanger.MarketOrder}, "<MKT 0x@0>"},
		{exchanger.Order{Type: exchanger.LimitOrder}, "<LIM 0x@0>"},
		{exchanger.Order{Type: exchanger.MarketOrder, Quantity: 10, Price: 99, Symbol: "TST"}, "<MKT 10xTST@99>"},
		{exchanger.Order{Type: exchanger.LimitOrder, Quantity: 10, Price: 99, Symbol: "TST"}, "<LIM 10xTST@99>"},
	}

	for _, tst := range tests {
		got := tst.input.String()
		if got != tst.exp {
			t.Errorf("%#v stringified to %q; expected %q", tst.input, got, tst.exp)
		}
	}
}

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
		{exchanger.Order{},
			"<SELL/MKT 0x@0>"},
		{exchanger.Order{OrderInfo: exchanger.Buy},
			"<BUY/MKT 0x@0>"},
		{exchanger.Order{OrderInfo: exchanger.Limit, Quantity: 10, Price: 99, Symbol: "TST"},
			"<SELL/LIM 10xTST@99>"},
		{exchanger.Order{OrderInfo: exchanger.Buy | exchanger.Limit, Quantity: 10, Price: 99, Symbol: "TST"},
			"<BUY/LIM 10xTST@99>"},
	}

	for _, tst := range tests {
		got := tst.input.String()
		if got != tst.exp {
			t.Errorf("%#v stringified to %q; expected %q", tst.input, got, tst.exp)
		}
	}
}

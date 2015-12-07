package exchanger_test

import (
	"github.com/djmitche/exchanger"
	"testing"
)

func TestTickString(t *testing.T) {
	tests := []struct {
		input exchanger.Tick
		exp   string
	}{
		{exchanger.Tick{}, "<??? 0x@0>"},
		{exchanger.Tick{Type: exchanger.QuoteTick, Quantity: 10, Price: 99, Symbol: "TST"}, "<Q 10xTST@99>"},
		{exchanger.Tick{Type: exchanger.ExecutionTick, Quantity: 10, Price: 99, Symbol: "TST"}, "<E 10xTST@99>"},
		{exchanger.Tick{Type: exchanger.CancellationTick, Quantity: 10, Price: 99, Symbol: "TST"}, "<C 10xTST@99>"},
	}

	for _, tst := range tests {
		got := tst.input.String()
		if got != tst.exp {
			t.Errorf("%#v stringified to %q; expected %q", tst.input, got, tst.exp)
		}
	}
}

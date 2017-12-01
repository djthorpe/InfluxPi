package influxpi_test

import (
	"testing"

	influx "github.com/djthorpe/InfluxPi"
)

func TestQuote_000(t *testing.T) {
	tests := map[string]string{
		"":                "",
		"a":               "a",
		"t1":              "t1",
		"_":               "_",
		"0":               "\"0\"",
		"a0":              "a0",
		"test string":     "\"test string\"",
		"test \"string\"": "\"test \\\"string\\\"\"",
		"test \\string\"": "\"test \\\\string\\\"\"",
	}
	for k, expected := range tests {
		actual := influx.QuoteIdentifier(k)
		if actual != expected {
			t.Errorf("For [%v], expected [%v], got [%v]", k, expected, actual)
		}
	}
}

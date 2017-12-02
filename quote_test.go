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

func TestUnquote_000(t *testing.T) {
	tests := map[string]string{
		"system,host=rpi3":          "system",
		"\"measurement\",host=rpi3": "measurement",
		"\"one two\",host=rpi3":     "one two",
		"one\\\"two,host=rpi3":      "one two",
		"\"one\\\"two\",host=rpi3":  "one two",
		"one\\,two,host=rpi3":       "one,two",
	}
	for k, expected := range tests {
		if actual, err := influx.UnquoteLine(k); err != nil {
			t.Errorf("For [%v], error %v", k, err)
		} else if actual != expected {
			t.Errorf("For [%v], expected [%v], got [%v]", k, expected, actual)
		}
	}
}

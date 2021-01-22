package collectors_test

import (
	"github.com/stretchr/testify/assert"
	"recipe-counter/collectors"
	"testing"
)

func TestDeliveryWindowCounter_Process(t *testing.T) {
	records := []*collectors.Record{
		{
			"Test 1",
			"10000",
			&collectors.Delivery{Weekday: "Monday", From: 11, To: 14},
		},
		{
			"Test 2",
			"10001",
			&collectors.Delivery{Weekday: "Tuesday", From: 12, To: 17},
		},
		{
			"Test 1",
			"10002",
			&collectors.Delivery{Weekday: "Wednesday", From: 9, To: 17},
		},
		{
			"Test 3",
			"10000",
			&collectors.Delivery{Weekday: "Friday", From: 8, To: 13},
		},
	}

	tests := map[string]struct {
		postcode, from, to string
		expected           int
	}{
		"Postcode is not provided": {
			postcode: "",
			from:     "1AM",
			to:       "1PM",
			expected: 0,
		},
		"Postcode provided, but is not within dataset": {
			postcode: "22222",
			from:     "1AM",
			to:       "1PM",
			expected: 0,
		},
		"Postcode provided, time window doesn't match": {
			postcode: "10000",
			from:     "5AM",
			to:       "12PM",
			expected: 0,
		},
		"Postcode provided, time window does match": {
			postcode: "10000",
			from:     "3AM",
			to:       "1PM",
			expected: 1,
		},
	}

	for name, data := range tests {
		c := collectors.NewDeliveryWindowCounter(data.postcode, data.from, data.to)
		for _, r := range records {
			c.Process(r)
		}

		report := c.Report()["count_per_postcode_and_time"].(map[string]interface{})
		assert.Equalf(t, data.expected, report["delivery_count"], name)
	}
}

package collectors_test

import (
	"github.com/stretchr/testify/assert"
	"recipe-counter/collectors"
	"testing"
)

func TestAmpmto24h(t *testing.T) {
	tests := []struct {
		timeStr  string
		expected int
		err      bool
	}{
		{timeStr: "1AM", expected: 1, err: false},
		{timeStr: "10AM", expected: 10, err: false},
		{timeStr: "12AM", expected: 0, err: false},
		{timeStr: "1PM", expected: 13, err: false},
		{timeStr: "10PM", expected: 22, err: false},
		{timeStr: "11PM", expected: 23, err: false},
		{timeStr: "12PM", expected: 12, err: false},
		{timeStr: "0AM", expected: 0, err: true},
		{timeStr: "13PM", expected: 0, err: true},
	}

	for _, data := range tests {
		result, err := collectors.Ampmto24h(data.timeStr)
		assert.Equal(t, data.expected, result)

		if data.err {
			assert.Error(t, err)
		} else {
			assert.Nil(t, err)
		}
	}
}

// ParseDelivery tests are simplified and doesn't test wrong message format as well as code doesn't handle it
func TestParseDelivery_panicWhenBadTimeFormat(t *testing.T) {
	assert.Panics(t, func() {
		collectors.ParseDelivery("Monday 28AM - 23PM")
	})

	assert.Panics(t, func() {
		collectors.ParseDelivery("Monday 0AM - 23PM")
	})

	assert.Panics(t, func() {
		collectors.ParseDelivery("Monday 1AM - 23PM")
	})
}

func TestParseDelivery(t *testing.T) {
	tests := []struct {
		delivery string
		expected collectors.Delivery
	}{
		{
			delivery: "Monday 1AM - 1PM",
			expected: collectors.Delivery{
				Weekday: "Monday",
				From:    1,
				To:      13,
			},
		},
		{
			delivery: "Wednesday 12AM - 1PM",
			expected: collectors.Delivery{
				Weekday: "Wednesday",
				From:    0,
				To:      13,
			},
		},
		{
			delivery: "Friday 8AM - 12PM",
			expected: collectors.Delivery{
				Weekday: "Friday",
				From:    8,
				To:      12,
			},
		},
	}

	for _, data := range tests {
		result := collectors.ParseDelivery(data.delivery)
		assert.Equal(t, data.expected.Weekday, result.Weekday)
		assert.Equal(t, data.expected.From, result.From)
		assert.Equal(t, data.expected.To, result.To)
	}
}

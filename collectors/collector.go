package collectors

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

var (
	deliveryRe *regexp.Regexp
	ampmRe     *regexp.Regexp
)

type Delivery struct {
	Weekday string
	From    int
	To      int
}

type Record struct {
	Recipe   string
	Postcode string
	Delivery *Delivery
}

type Collector interface {
	Process(r *Record)
	Report() map[string]interface{}
}

func init() {
	deliveryRe = regexp.MustCompile("(?i)([a-z]+) ([0-9a-z]+)+[- ]+([0-9a-z]+)")
	ampmRe = regexp.MustCompile("(?i)([0-9]+)([amp]{2})")
}

func Ampmto24h(t string) int {
	parts := ampmRe.FindStringSubmatch(t)
	if parts[1] == "" || parts[2] == "" {
		panic(fmt.Errorf("invalid input: %v, 00AM or 00PM expected", t))
	}

	hours, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		panic(err)
	}

	hours = math.Mod(hours, 12)

	if strings.ToUpper(parts[2]) == "PM" {
		hours += 12
	}

	return int(hours)
}

// We assume now we have only one format as defined int the task
// Also I'll skip error handling for values that didn't match the pattern as tasks says it should
func ParseDelivery(val string) *Delivery {
	parts := deliveryRe.FindStringSubmatch(val)

	return &Delivery{
		Weekday: parts[1],
		From:    Ampmto24h(parts[2]),
		To:      Ampmto24h(parts[3]),
	}
}

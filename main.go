package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"recipe-counter/collectors"
	"recipe-counter/parser"
	"strings"
)

type searchPatternVar []string

func (s searchPatternVar) String() string {
	return strings.Join(s, " ")
}

func (s *searchPatternVar) Set(val string) error {
	*s = append(*s, val)

	return nil
}

var (
	searchPostcode   string
	searchWindowFrom string
	searchWindowTo   string
	recipePattern    searchPatternVar
)

func init() {
	flag.StringVar(&searchPostcode, "postcode", "", "Search for a given postcode")
	flag.StringVar(&searchWindowFrom, "window-from", "12AM", "Search for a given postcode in window from")
	flag.StringVar(&searchWindowTo, "window-to", "11PM", "Search for a given postcode in window to")
	flag.Var(&recipePattern, "recipe", "Search for a recipes containing given pattern, can be specified multiple times")
}

func main() {
	flag.Parse()

	inputSrc := flag.Arg(0)
	p := parser.NewStreamParser(inputSrc)
	stats := collectors.NewStatsAggregate(
		[]collectors.Collector{
			collectors.NewRecipesCollector(recipePattern),
			collectors.NewBusiestPostcodeCollector(),
			collectors.NewDeliveryWindowCounter(searchPostcode, searchWindowFrom, searchWindowTo),
		}...,
	)

	p.Read(func(r *parser.Record) {
		// Ideally we would have mappers layer to map parser Record to collectors record
		// But I'll skip it for sake of simplicity here
		stats.Process(&collectors.Record{
			Recipe:   r.Recipe,
			Postcode: r.Postcode,
			Delivery: collectors.ParseDelivery(r.Delivery),
		})
	})

	report, err := json.MarshalIndent(stats.Report(), "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Print(string(report))
}

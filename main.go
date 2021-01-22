package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
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
	inputSrc         string
	searchPostcode   string
	searchWindowFrom string
	searchWindowTo   string
	recipePattern    searchPatternVar
)

func init() {
	flag.StringVar(&inputSrc, "src", "", "Source json file path with deliveries data to collect")
	flag.StringVar(&searchPostcode, "postcode", "", "Search for a given postcode")
	flag.StringVar(&searchWindowFrom, "window-from", "12AM", "Search for a given postcode in window from")
	flag.StringVar(&searchWindowTo, "window-to", "11PM", "Search for a given postcode in window to")
	flag.Var(&recipePattern, "recipe", "Search for a recipes containing given pattern, can be specified multiple times")
}

func main() {
	flag.Parse()
	defer func() {
		if err := recover(); err != nil {
			log.Fatal(err)
		}
	}()

	if inputSrc == "" {
		fmt.Printf("Source json file is not provided, please see usage for available options:\n")
		flag.Usage()
		os.Exit(0)
	}

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

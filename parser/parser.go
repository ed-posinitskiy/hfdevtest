package parser

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Record struct {
	Postcode, Recipe, Delivery string
}

type StreamParser struct {
	reader io.Reader
}

func (p StreamParser) Read(listener func(r *Record)) {
	decoder := json.NewDecoder(p.reader)

	// Verify if given input is an array
	t, err := decoder.Token()
	if err != nil {
		panic(err)
	}

	if fmt.Sprintf("%v", t) != "[" {
		panic(fmt.Errorf("read error: expected token [, recieved %v", t))
	}

	for decoder.More() {
		var r Record
		if err := decoder.Decode(&r); err != nil {
			panic(err)
		}
		listener(&r)
	}
}

func NewStreamParser(src string) *StreamParser {
	file, err := os.Open(src)
	if err != nil {
		panic(err)
	}

	return &StreamParser{file}
}

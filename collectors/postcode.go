package collectors

type BusiestPostcodeCollector struct {
	busiest   string
	postcodes map[string]int
}

func (c *BusiestPostcodeCollector) Process(r *Record) {
	if _, ok := c.postcodes[r.Postcode]; !ok {
		c.postcodes[r.Postcode] = 0
	}

	c.postcodes[r.Postcode] += 1

	if c.busiest == "" {
		c.busiest = r.Postcode
		return
	}

	if c.busiest == r.Postcode {
		return
	}

	if c.postcodes[r.Postcode] > c.postcodes[c.busiest] {
		c.busiest = r.Postcode
	}
}

func (c BusiestPostcodeCollector) Report() map[string]interface{} {
	return map[string]interface{}{
		"busiest_postcode": map[string]interface{}{
			"postcode":       c.busiest,
			"delivery_count": c.postcodes[c.busiest],
		},
	}
}

func NewBusiestPostcodeCollector() *BusiestPostcodeCollector {
	return &BusiestPostcodeCollector{
		busiest:   "",
		postcodes: make(map[string]int),
	}
}

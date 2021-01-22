package collectors

type DeliveryWindowCounter struct {
	from     int
	fromAm   string
	to       int
	toPm     string
	postcode string
	counter  int
}

func (c *DeliveryWindowCounter) Process(r *Record) {
	if r.Postcode != c.postcode {
		return
	}

	if r.Delivery.From < c.from || r.Delivery.To > c.to {
		return
	}

	c.counter += 1
}

func (c DeliveryWindowCounter) Report() map[string]interface{} {
	return map[string]interface{}{
		"count_per_postcode_and_time": map[string]interface{}{
			"postcode":       c.postcode,
			"from":           c.fromAm,
			"to":             c.toPm,
			"delivery_count": c.counter,
		},
	}
}

func NewDeliveryWindowCounter(postcode, from, to string) *DeliveryWindowCounter {
	var (
		fromTime, toTime int
		err              error
	)

	fromTime, err = Ampmto24h(from)
	if err != nil {
		panic(err)
	}

	toTime, err = Ampmto24h(to)
	if err != nil {
		panic(err)
	}

	return &DeliveryWindowCounter{
		from:     fromTime,
		fromAm:   from,
		to:       toTime,
		toPm:     to,
		postcode: postcode,
		counter:  0,
	}
}

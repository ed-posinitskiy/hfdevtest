package collectors

type StatsAggregate struct {
	collectors []Collector
}

func (s StatsAggregate) Process(r *Record) {
	for _, c := range s.collectors {
		c.Process(r)
	}
}

func (s StatsAggregate) Report() map[string]interface{} {
	report := make(map[string]interface{})
	for _, c := range s.collectors {
		// Ideally it should be recursive merge, but I'll skip it for simplicity
		for k, v := range c.Report() {
			report[k] = v
		}
	}

	return report
}

func NewStatsAggregate(collectors ...Collector) *StatsAggregate {
	return &StatsAggregate{
		collectors: collectors,
	}
}

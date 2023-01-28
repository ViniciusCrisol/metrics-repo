package metric

type Repo interface {
	// Create persists a metric
	Create(m *Metric) error
}

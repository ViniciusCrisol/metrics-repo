package metric

type Repo interface {
	Create(m *Metric) error
}

package metric

import (
	"encoding/json"

	"github.com/ViniciusCrisol/metrics-repo/insert-data-service/pkg/input"
)

type Metric struct {
	Data    string `json:"data"`
	AppName string `json:"app_name"`
}

func FromInput(i *input.Input) (*Metric, error) {
	m := &Metric{}

	if err := json.Unmarshal([]byte(i.Data), m); err != nil {
		// TODO: Log it!
		return nil, err
	}
	return m, nil
}

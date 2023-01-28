package metric

import (
	"encoding/json"

	"github.com/ViniciusCrisol/metrics-repo/insert-data-service/log"
	"github.com/ViniciusCrisol/metrics-repo/insert-data-service/pkg/input"
)

type Metric struct {
	Data    string `json:"data"`
	AppName string `json:"app_name"`
}

func FromInput(i *input.Input) (*Metric, error) {
	m := &Metric{}

	if err := json.Unmarshal([]byte(i.Data), m); err != nil {
		log.Logger.Error(
			"Failed to unmarshal input",
			log.Error(err),
			log.String("input_id", i.ID),
			log.String("message_id", i.ID),
			log.String("input_data", i.Data),
		)
		return nil, err
	}
	return m, nil
}

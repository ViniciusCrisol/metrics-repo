package input

type Input struct {
	ID   string
	Data string
}

func NewInput(id, data string) *Input {
	return &Input{
		ID:   id,
		Data: data,
	}
}

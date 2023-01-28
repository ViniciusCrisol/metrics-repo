package input

type Repo interface {
	Get() ([]*Input, error)

	Delete(i *Input) error
}

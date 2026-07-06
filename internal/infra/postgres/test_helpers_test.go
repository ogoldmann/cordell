package postgres_test

type fixedIDGenerator struct {
	id string
}

func (g fixedIDGenerator) NewID() (string, error) {
	return g.id, nil
}

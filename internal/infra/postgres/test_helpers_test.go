package postgres_test

type fixedIDGenerator struct {
	id string
}

func (g fixedIDGenerator) NewID() string {
	return g.id
}

package mock

type IDGenerator struct {
	Generates int64
}

func (i *IDGenerator) Generate() int64 {
	return i.Generates
}

package mock

type CardIDGenerator struct {
	Generates int64
}

func (i *CardIDGenerator) Generate() int64 {
	return i.Generates
}

type StringIDGenerator struct {
	MockID string
}

func (i *StringIDGenerator) Generate() string {
	return i.MockID
}

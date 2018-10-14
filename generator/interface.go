package generator

type CardIDGenerator interface {
	Generate() int64
}

type StringIDGenerator interface {
	Generate() string
}

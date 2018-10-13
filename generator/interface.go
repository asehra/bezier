package generator

type IDGenerator interface {
	Generate() int64
}

package generator

type CardIDGenerator interface {
	Generate() int64
}

type MerchantIDGenerator interface {
	Generate() string
}

type StringIDGenerator interface {
	Generate() string
}

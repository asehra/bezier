package generator

type CardIDGenerator interface {
	Generate() int64
}

type MerchantIDGenerator interface {
	Generate() string
}

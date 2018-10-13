package mock

type CardIDGenerator struct {
	Generates int64
}

func (i *CardIDGenerator) Generate() int64 {
	return i.Generates
}

type MerchantIDGenerator struct {
	Generates string
}

func (i *MerchantIDGenerator) Generate() string {
	return i.Generates
}

package generator

import "strconv"

type CardIDIncrementor struct {
	LastID int64
}

func (c *CardIDIncrementor) Generate() int64 {
	c.LastID = c.LastID + 1
	return c.LastID
}

type MerchantIDIncrementor struct {
	Prefix string
	LastID int
}

func (c *MerchantIDIncrementor) Generate() string {
	c.LastID = c.LastID + 1
	return c.Prefix + strconv.Itoa(c.LastID)
}

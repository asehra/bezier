package generator

type CardID struct {
	LastID int64
}

func (c *CardID) Generate() int64 {
	c.LastID = c.LastID + 1
	return c.LastID
}

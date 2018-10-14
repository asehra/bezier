package generator

import "strconv"

type CardIDIncrementor struct {
	LastID int64
}

func (c *CardIDIncrementor) Generate() int64 {
	c.LastID = c.LastID + 1
	return c.LastID
}

type StringIDIncrementor struct {
	Prefix string
	LastID int
}

func (i *StringIDIncrementor) Generate() string {
	i.LastID = i.LastID + 1
	return i.Prefix + strconv.Itoa(i.LastID)
}

package bezier

type Card struct {
	Number           int64
	AvailableBalance int32
}

func CreateCard() *Card {
	return &Card{
		Number: 9000000000000001,
	}
}

package model

type Card struct {
	Number           int64 `json:"card_number"`
	AvailableBalance int32 `json:"available_balance"`
}

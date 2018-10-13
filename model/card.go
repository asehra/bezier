package model

type Card struct {
	Number           int64 `json:"card_number"`
	AvailableBalance int32 `json:"available_balance"`
	BlockedBalance   int32 `json:"blocked_balance"`
	TotalLoaded      int32 `json:"total_loaded"`
}

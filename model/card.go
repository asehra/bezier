package model

type Card struct {
	Number           int64 `json:"card_number"`
	AvailableBalance int   `json:"available_balance"`
	BlockedBalance   int   `json:"blocked_balance"`
	TotalLoaded      int   `json:"total_loaded"`
}

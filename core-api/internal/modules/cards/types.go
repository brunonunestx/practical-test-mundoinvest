package cards

type CardUpdateDto struct {
	EventID    string `json:"event_id" validate:"required"`
	CardID     string `json:"card_id" validate:"required"`
	ClientMail string `json:"cliente_email" validate:"required,email"`
	Timestamp  string `json:"timestamp" validate:"required"`
}

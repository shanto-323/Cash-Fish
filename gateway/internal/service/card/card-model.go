package card

type CardsResponseMetadata struct {
	UID         string `json:"uid"`
	Number      string `json:"number"`
	Brand       string `json:"brand"`
	ExpiryMonth int    `json:"expiry_month"`
	ExpiryYear  int    `json:"expiry_year"`
}

type CardID struct {
	ID string `json:"id"` // CARD ID
}

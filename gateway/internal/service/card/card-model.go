package card

type CardsResponseMetadata struct {
	ID          string `json:"id"`
	Number      string `json:"number"`
	Brand       string `json:"brand"`
	ExpiryMonth int    `json:"expiry_month"`
	ExpiryYear  int    `json:"expiry_year"`
}

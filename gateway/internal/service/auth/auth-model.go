package auth

type UserResponseModel struct {
	ID       string                  `json:"id"`
	Username string                  `json:"username"`
	Password string                  `json:"password"`
	Email    string                  `json:"email"`
	Cards    []CardsResponseMetadata `json:"cards"`
	Token    TokenMetadata           `json:"token"`
}

type TokenMetadata struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

type CardsResponseMetadata struct {
	ID          string `json:"id"`
	Number      string `json:"number"`
	Brand       string `json:"brand"`
	ExpiryMonth int    `json:"expiry_month"`
	ExpiryYear  int    `json:"expiry_year"`
}

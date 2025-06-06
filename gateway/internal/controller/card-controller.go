package controller

import (
	"gateway/internal/service/card"

	"github.com/gorilla/mux"
)

type CardController struct {
	router     *mux.Router
	cardClient *card.CardClient
}

func NewCardController(router *mux.Router, cardClient *card.CardClient) *CardController {
	router.PathPrefix("/card").Subrouter()
	return &CardController{
		router:     router,
		cardClient: cardClient,
	}
}

func (c *CardController) RegisterRoutes() {
}

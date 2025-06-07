package controller

import (
	"context"
	"encoding/json"
	"gateway/internal/service/card"
	"gateway/pkg"
	"net/http"
	"time"

	"gateway/internal/middlerware"

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
	c.router.Use(middlerware.JwtMiddleware)
	c.router.HandleFunc("/add", middlerware.HandleFunc(c.AddCard)).Methods("POST")
	c.router.HandleFunc("/all/{uid}", middlerware.HandleFunc(c.GetAllCard)).Methods("GET")
	c.router.HandleFunc("/all/{uid}", middlerware.HandleFunc(c.DeleteAllCard)).Methods("DELETE")
	c.router.HandleFunc("/remove/{uid}/{id}", middlerware.HandleFunc(c.RemoveCard)).Methods("DELETE")
}

func (c *CardController) AddCard(w http.ResponseWriter, r *http.Request) error {
	card := card.CardsResponseMetadata{}
	if err := json.NewDecoder(r.Body).Decode(&card); err != nil {
		return err
	}
	defer r.Body.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, err := c.cardClient.CardClientAddCard(ctx, card.UID, card.Number, card.Brand, int32(card.ExpiryMonth), int32(card.ExpiryYear))
	if err != nil {
		return err
	}

	return pkg.WriteJson(w, http.StatusOK, resp)
}

func (c *CardController) GetAllCard(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	uid := vars["uid"]

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, err := c.cardClient.CardClientGetALlCards(ctx, uid)
	if err != nil {
		return err
	}

	return pkg.WriteJson(w, http.StatusOK, resp)
}

func (c *CardController) DeleteAllCard(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	uid := vars["uid"]

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, err := c.cardClient.CardClientDeleteCards(ctx, uid)
	if err != nil {
		return err
	}

	return pkg.WriteJson(w, http.StatusOK, resp)
}

func (c *CardController) RemoveCard(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	uid := vars["uid"]
	id := vars["id"]

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, err := c.cardClient.CardClientRemoveCard(ctx, uid, id)
	if err != nil {
		return err
	}

	return pkg.WriteJson(w, http.StatusOK, resp)
}

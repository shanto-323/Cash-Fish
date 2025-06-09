package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"gateway/internal/service/card"
	"gateway/pkg"

	"gateway/internal/middlerware"

	"github.com/gorilla/mux"
)

type CardController struct {
	router     *mux.Router
	cardClient *card.CardClient
}

func NewCardController(router *mux.Router, cardClient *card.CardClient) *CardController {
	newCardRouter := router.PathPrefix("/card").Subrouter()
	return &CardController{
		router:     newCardRouter,
		cardClient: cardClient,
	}
}

func (c *CardController) RegisterRoutes() {
	c.router.Use(middlerware.JwtMiddleware)
	c.router.HandleFunc("/add", middlerware.HandleFunc(c.AddCard)).Methods("POST") // ID -> USER_ID
	c.router.HandleFunc("/all", middlerware.HandleFunc(c.GetAllCard)).Methods("GET")
	c.router.HandleFunc("/all", middlerware.HandleFunc(c.DeleteAllCard)).Methods("DELETE")
	c.router.HandleFunc("/remove/{cid}", middlerware.HandleFunc(c.RemoveCard)).Methods("DELETE") // CID -> CARD_ID
}

func (c *CardController) AddCard(w http.ResponseWriter, r *http.Request) error {
	card := &card.CardsResponseMetadata{}
	if err := json.NewDecoder(r.Body).Decode(&card); err != nil {
		return fmt.Errorf("json marshaling error%s", err)
	}
	defer r.Body.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, err := c.cardClient.CardClientAddCard(ctx, card.UID, card.Number, card.Brand, int32(card.ExpiryMonth), int32(card.ExpiryYear))
	if err != nil {
		return fmt.Errorf("card client error%s", err)
	}

	return pkg.WriteJson(w, http.StatusOK, resp)
}

func (c *CardController) GetAllCard(w http.ResponseWriter, r *http.Request) error {
	uid := r.URL.Query().Get("id")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, err := c.cardClient.CardClientGetALlCards(ctx, uid)
	if err != nil {
		return err
	}

	return pkg.WriteJson(w, http.StatusOK, resp)
}

func (c *CardController) DeleteAllCard(w http.ResponseWriter, r *http.Request) error {
	uid := r.URL.Query().Get("id")
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
	id := r.URL.Query().Get("id")
	cid := vars["cid"]

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, err := c.cardClient.CardClientRemoveCard(ctx, id, cid)
	if err != nil {
		return err
	}

	return pkg.WriteJson(w, http.StatusOK, resp)
}

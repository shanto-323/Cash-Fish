package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"gateway/internal/middlerware"
	"gateway/internal/service/auth"
	"gateway/pkg"

	"github.com/gorilla/mux"
)

type AuthController struct {
	router     *mux.Router
	authClient *auth.AuthClient
}

func NewAuthController(router *mux.Router, authClient *auth.AuthClient) *AuthController {
	authRouter := router.PathPrefix("/auth").Subrouter()
	return &AuthController{
		router:     authRouter,
		authClient: authClient,
	}
}

func (c *AuthController) RegisterRoutes() {
	c.router.HandleFunc("/signup", middlerware.HandleFunc(c.SignUp)).Methods("POST")
	c.router.HandleFunc("/signin", middlerware.HandleFunc(c.SignIn)).Methods("POST")
	c.router.HandleFunc("/token", middlerware.HandleFunc(c.NewToken)).Methods("GET")

	protected := c.router.PathPrefix("/user").Subrouter()
	protected.Use(middlerware.JwtMiddleware)
	protected.HandleFunc("/signout", middlerware.HandleFunc(c.SignOut)).Methods("POST")
	protected.HandleFunc("/delete", middlerware.HandleFunc(c.DeleteUser)).Methods("DELETE")
}

func (c *AuthController) SignUp(w http.ResponseWriter, r *http.Request) error {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	user := auth.UserModel{}
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		return err
	}

	resp, err := c.authClient.AuthClientSignUP(ctx, user.Username, user.Email, user.Password)
	if err != nil {
		return err
	}
	http.SetCookie(
		w,
		&http.Cookie{
			Name:  "access_token",
			Value: resp.Token.Token,
			Path:  "/api/v1/",
		},
	)

	return pkg.WriteJson(w, http.StatusOK, resp)
}

func (c *AuthController) SignIn(w http.ResponseWriter, r *http.Request) error {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	user := auth.UserModel{}
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		return err
	}

	resp, err := c.authClient.AuthClientSignIN(ctx, user.Email, user.Password)
	if err != nil {
		return err
	}

	http.SetCookie(
		w,
		&http.Cookie{
			Name:  "access_token",
			Value: resp.Token.Token,
			Path:  "/api/v1/",
		},
	)
	return pkg.WriteJson(w, http.StatusOK, resp)
}

func (c *AuthController) SignOut(w http.ResponseWriter, r *http.Request) error {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	id := r.URL.Query().Get("id")
	resp, err := c.authClient.AuthClientSignOut(ctx, id)
	if err != nil {
		return fmt.Errorf("signout error:%s", err)
	}
	return pkg.WriteJson(w, http.StatusOK, *resp)
}

func (c *AuthController) DeleteUser(w http.ResponseWriter, r *http.Request) error {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	id := r.URL.Query().Get("id")
	resp, err := c.authClient.AuthClientDeleteUser(ctx, id)
	if err != nil {
		return fmt.Errorf("signout error:%s", err)
	}
	return pkg.WriteJson(w, http.StatusOK, *resp)
}

func (c *AuthController) NewToken(w http.ResponseWriter, r *http.Request) error {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	id := r.URL.Query().Get("id")
	token := auth.Token{}
	if err := json.NewDecoder(r.Body).Decode(&token); err != nil {
		return err
	}

	resp, err := c.authClient.AuthClientNewToken(ctx, id, token.RefreshToken)
	if err != nil {
		return err
	}

	http.SetCookie(
		w,
		&http.Cookie{
			Name:  "access_token",
			Value: *resp,
			Path:  "./",
		},
	)
	return pkg.WriteJson(w, http.StatusOK, *resp)
}

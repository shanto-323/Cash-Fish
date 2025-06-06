package pkg

import (
	"encoding/json"
	"net/http"
)

func WriteJson(w http.ResponseWriter, status int, msg any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	type message struct {
		Statue int `json:"code"`
		Msg    any `json:"msg"`
	}

	return json.NewEncoder(w).Encode(message{
		Statue: status,
		Msg:    msg,
	})
}

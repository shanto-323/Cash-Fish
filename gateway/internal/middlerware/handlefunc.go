package middlerware

import (
	"net/http"

	"gateway/pkg"
)

func HandleFunc(f func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			pkg.WriteJson(w, http.StatusBadRequest, err.Error())
			return
		}
	}
}

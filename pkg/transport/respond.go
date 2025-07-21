package transport

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Err string `json:"error"`
}

func Http(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func RespondWithError(w http.ResponseWriter, code int, err error, msg string) {
	var ermsg string

	if msg != "" {
		ermsg = msg + ": " + err.Error()
	} else {
		ermsg = err.Error()
	}

	http.Error(w, ermsg, code)
}

func ErrorRespond(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	response := ErrorResponse{
		Err: msg,
	}

	_ = json.NewEncoder(w).Encode(response)
}

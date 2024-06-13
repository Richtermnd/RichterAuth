package utils

import (
	"encoding/json"
	"net/http"

	"github.com/Richtermnd/RichterAuth/internal/errs"
)

type httpCoder interface {
	HttpCode() int
}

// Encode Encode response to json or internal error on failed marshalling
func Encode(w http.ResponseWriter, r *http.Request, v interface{}) error {
	if errCoder, ok := v.(httpCoder); ok {
		w.WriteHeader(errCoder.HttpCode())
	}
	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		json.NewEncoder(w).Encode(errs.ErrInternal(err))
	}
	return nil
}

// Decode Decode request from json to v
//
// v must be a pointer
func Decode(w http.ResponseWriter, r *http.Request, v interface{}) error {
	err := json.NewDecoder(r.Body).Decode(v)
	if err != nil {
	}
	return err
}

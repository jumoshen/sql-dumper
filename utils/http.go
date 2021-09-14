package utils

import (
	"encoding/json"
	"net/http"

	"dumper/constant"
)

const (
	// ApplicationJson means application/json.
	ApplicationJson = "application/json"
	// ContentEncoding means Content-Encoding.
	ContentEncoding = "Content-Encoding"
	// ContentSecurity means X-Content-Security.
	ContentSecurity = "X-Content-Security"
	// ContentType means Content-Type.
	ContentType = "Content-Type"
	// KeyField means key.
	KeyField = "key"
	// SecretField means secret.
	SecretField = "secret"
	// TypeField means type.
	TypeField = "type"
	// CryptionType means cryption.
	CryptionType = 1
)

// OkJson writes v into w with 200 OK.
func OkJson(w http.ResponseWriter, v interface{}) {
	WriteJson(w, http.StatusOK, v)
}

// WriteJson writes v as json string into w with code.
func WriteJson(w http.ResponseWriter, code int, v interface{}) {
	w.Header().Set(ContentType, ApplicationJson)
	w.WriteHeader(code)

	if bs, err := json.Marshal(v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else if n, err := w.Write(bs); err != nil {
		// http.ErrHandlerTimeout has been handled by http.TimeoutHandler,
		// so it's ignored here.
		if err != http.ErrHandlerTimeout {
			constant.Logger.Errorf("write response failed, error: %s", err)
		}
	} else if n < len(bs) {
		constant.Logger.Errorf("actual bytes: %d, written bytes: %d", len(bs), n)
	}
}

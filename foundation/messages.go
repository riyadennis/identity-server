package foundation

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/http"
)

type Response struct {
	Status    int    `json:"status"`
	Message   string `json:"message"`
	ErrorCode string `json:"error-code"`
}

func ErrorResponse(w http.ResponseWriter, code int, errr error, customCode string) {
	w.Header().Set("Content-Type", "application/json")
	err := JSONResponse(w, code, errr.Error(), customCode)
	if err != nil {
		logrus.Error(err)
	}
}

func JSONResponse(w http.ResponseWriter, status int,
	message, errCode string) error {
	w.WriteHeader(status)
	res := NewResponse(status, message, errCode)
	err := json.NewEncoder(w).Encode(res)
	if err != nil {
		return err
	}
	return nil
}

func NewResponse(status int, message, errCode string) *Response {
	return &Response{
		Status:    status,
		Message:   message,
		ErrorCode: errCode,
	}
}

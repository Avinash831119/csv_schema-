package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (e *MyError) WriteToResponse(w http.ResponseWriter) {
	w.WriteHeader(e.HTTPStatus)
	fmt.Fprintf(w, e.ToJSON())
}

func (e *MyError) ToJSON() string {
	j, err := json.Marshal(e)
	if err != nil {
		return `{"code":50099,"message":"ScrapError.JSONStr: json.Marshal() failed"}`
	}
	return string(j)
}

type MyError struct {
	HTTPStatus int    `json:"-"`
	Code       int    `json:"code"`
	Message    string `json:"message"`
}

func (e *MyError) Error() string {
	return fmt.Sprintf("HTTPStatus: %v, Code: %v, Message: %q",
		e.HTTPStatus, e.Code, e.Message)
}

func JSONMarshalErr(message string, code int) *MyError {
	return &MyError{
		HTTPStatus: http.StatusInternalServerError,
		Code:       code,
		Message:    message,
	}
}

func JSONMarshalBadRequest(message string, code int) *MyError {
	return &MyError{
		HTTPStatus: http.StatusBadRequest,
		Code:       code,
		Message:    message,
	}
}

func JSONMarshalNoContent(message string, code int) *MyError {
	return &MyError{
		HTTPStatus: http.StatusNoContent,
		Code:       code,
		Message:    message,
	}
}

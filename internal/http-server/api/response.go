package api

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

// func WriteRsponse()

const (
	StatusOk    = "Ok"
	StatusError = "Error"
)

type CommonResponse struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

func ResponseError(err string) CommonResponse {
	return CommonResponse{
		Status: StatusError,
		Error:  err,
	}
}

func ValidationError(errs validator.ValidationErrors) CommonResponse {
	var errMsgs []string

	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is a required field", err.Field()))
		case "url":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not valid Url", err.Field()))
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not valid", err.Field()))
		}
	}

	return CommonResponse{
		Status: StatusError,
		Error:  strings.Join(errMsgs, ", "),
	}
}

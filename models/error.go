package models

import (
	"fmt"
)

type ApiError struct {
	isError    bool   `json:"error"`
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

func NewApiError(statusCode int, format string, a ...interface{}) ApiError {
	if statusCode <= 0 {
		statusCode = 400
	}

	return ApiError{
		isError:    true,
		StatusCode: statusCode,
		Message:    fmt.Sprintf(format, a...),
	}
}

func (err ApiError) Error() string {
	return fmt.Sprintf("HTTP %v - %v", err.StatusCode, err.Message)
}

package v1

import (
	"fmt"
)

type ApiError struct {
	isError    bool   `json:"error"`
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

func NewApiError(statusCode int, message string) ApiError {
	if statusCode <= 0 {
		statusCode = 400
	}

	return ApiError{
		isError:    true,
		StatusCode: statusCode,
		Message:    message,
	}
}

func (err ApiError) Error() string {
	return fmt.Sprintf("HTTP %v - %v", err.StatusCode, err.Message)
}

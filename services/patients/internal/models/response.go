package models

type JsonResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorJson  `json:"error,omitempty"`
}

type ErrorJson struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}


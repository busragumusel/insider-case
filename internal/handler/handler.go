package handler

type APIResult struct {
	Code    string      `json:"code,omitempty"`
	Message string      `json:"message,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
	Data    interface{} `json:"data"`
}

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

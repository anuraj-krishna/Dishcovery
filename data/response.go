package data

type SuccessResponse struct {
	Status     string `json:"status"`
	StatusCode int    `json:"status_code"`
	Count      int    `json:"count,omitempty"`
	Data       any    `json:"data,omitempty"`
}

type FailureResponse struct {
	Status       string `json:"status"`
	StatusCode   int    `json:"status_code"`
	Error        string `json:"error,omitempty"`
	ErrorDetails string `json:"details,omitempty"`
	Data         any    `json:"data,omitempty"`
}

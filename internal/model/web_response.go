package model

// WebResponse is the standard API envelope.
type WebResponse[T any] struct {
	Data T `json:"data"`
}

// WebResponsePaged wraps paginated results.
type WebResponsePaged[T any] struct {
	Data       []T        `json:"data"`
	Pagination Pagination `json:"pagination"`
}

type Pagination struct {
	NextCursor string `json:"next_cursor,omitempty"`
	HasMore    bool   `json:"has_more"`
	Total      int64  `json:"total,omitempty"`
}

// ErrorResponse is returned on all 4xx/5xx responses.
type ErrorResponse struct {
	Errors []ErrorItem `json:"errors"`
}

type ErrorItem struct {
	Field   string `json:"field,omitempty"`
	Message string `json:"message"`
}

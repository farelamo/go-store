package schema

type Response struct {
	Page     int    `json:"page,omitempty"`
	PageSize int    `json:"page_size,omitempty"`
	Total    int64  `json:"total,omitempty"`
	Message  string `json:"message"`
	Data     any    `json:"data"`
}

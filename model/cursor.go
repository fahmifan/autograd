package model

// Cursor ..
type Cursor struct {
	Limit  int64       `json:"limit"`
	Page   int64       `json:"page"`
	Offset int64       `json:"offset"`
	Sort   string      `json:"sort"`
	Rows   interface{} `json:"rows"`
}

// CursorRequest ..
type CursorRequest struct {
	Limit int64
	Page  int64
	Sort  string
}

// CursorResponse ..
type CursorResponse struct {
	Limit int64       `json:"limit"`
	Page  int64       `json:"page"`
	Sort  string      `json:"sort"`
	Rows  interface{} `json:"rows"`
}

package model

// Pagination ..
type Pagination struct {
	Limit  int64       `json:"limit"`
	Page   int64       `json:"page"`
	Offset int64       `json:"offset"`
	Sort   string      `json:"sort"`
	Rows   interface{} `json:"rows"`
}

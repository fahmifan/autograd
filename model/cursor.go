package model

// Cursor ..
type Cursor struct {
	Size   int64       `json:"size"`
	Page   int64       `json:"page"`
	Offset int64       `json:"offset"`
	Sort   string      `json:"sort"`
	Data   interface{} `json:"data"`
}

// CursorRequest ..
type CursorRequest struct {
	Size int64
	Page int64
	Sort string
}

// CursorResponse ..
type CursorResponse struct {
	Size int64       `json:"size"`
	Page int64       `json:"page"`
	Sort string      `json:"sort"`
	Data interface{} `json:"data"`
}

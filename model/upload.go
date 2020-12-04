package model

// Upload ..
type Upload struct {
	FileURL    string
	SourceCode string
}

// UploadRequest ..
type UploadRequest struct {
	SourceCode string `json:"sourceCode"`
}

// UploadResponse ..
type UploadResponse struct {
	FileURL string `json:"fileURL"`
}

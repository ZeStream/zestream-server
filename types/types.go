package types

// GenerateURLRequest is a struct that represents request body of /generate-presigned-url endpoint
type GenerateURLRequest struct {
	FileName string `json:"fileName"`
}

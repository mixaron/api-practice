package dto

type AttachmentResponse struct {
	URL string `json:"url"`
}

type AttachmentRequest struct {
	FileName string `json:"fileName"`
	URL      string `json:"url"`
}

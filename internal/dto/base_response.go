package dto

import "time"

type BaseResponse[T any] struct {
	Status  string `json:"status" example:"success"`
	Message string `json:"message" example:"operation done"`
	Data    *T     `json:"data,omitempty"`
}

// TokenBaseResponse, ErrorBaseResponse, SuccessRegisterBaseResponse для документации
type TokenBaseResponse struct {
	Status  string        `json:"status" example:"success"`
	Message string        `json:"message" example:"Authenticated"`
	Data    TokenResponse `json:"data"`
}

type ErrorBaseResponse struct {
	Status  string `json:"status" example:"error"`
	Message string `json:"message" example:"Error message"`
}

type SuccessBaseResponse struct {
	Status  string `json:"status" example:"success"`
	Message string `json:"message" example:"Success message"`
}

type SuccessProfileBaseResponse struct {
	Status  string          `json:"status" example:"success"`
	Message string          `json:"message" example:"Success message"`
	Data    ProfileResponse `json:"data"`
}

type ArticleBaseResponse struct {
	Title       string                   `json:"title" example:"title example"`
	Content     string                   `json:"content" example:"content example"`
	PreviewURL  string                   `json:"preview_url" example:"https://example.com/images/article1.jpg"`
	IsPublished bool                     `json:"is_published" example:"true"`
	AuthorID    uint                     `json:"author_id" example:"1"`
	CreatedAt   time.Time                `json:"created_at" example:"2025-05-31T10:30:00Z"`
	Attachments []AttachmentBaseResponse `json:"attachments"`
}

type AttachmentBaseResponse struct {
	URL      string `json:"url" example:"https://example.com/files/document.pdf"`
	FileName string `json:"file_name" example:"document.pdf"`
	FileSize int64  `json:"file_size" example:"1024"`
}

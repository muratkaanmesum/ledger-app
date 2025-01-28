package dtos

type PaginationRequest struct {
	Page  uint `json:"page"`
	Count uint `json:"count"`
}

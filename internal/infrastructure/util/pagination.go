package util

type Pagination struct {
	TotalCount int64 `json:"total_count"`
	TotalPages int64 `json:"total_pages"`
	Page       int   `json:"page"`
	Size       int   `json:"size"`
}

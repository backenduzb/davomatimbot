package crud

type PaginatedResponse[T any] struct {
	Count int64 `json:"count"`
	TotalPages int64 `json:"total_pages"`
	CurrentPage int64 `json:"current_page"`
	PageSize int64 	`json:"page_size"`
	Results []T `json:"results"`
}

type PaginationParams struct {
	Page int `form:"page" binding:"omitempty,gte=1"`
	PageSize int `form:"page_size" binding:"omitempty,gte=1,lte=100"`
}
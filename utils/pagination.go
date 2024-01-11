package utils

type DataTableResponse struct {
	Pagination Pagination  `json:"pagination"`
	Sort       string      `json:"sort"`
	Data       interface{} `json:"data"`
}
type Pagination struct {
	CurrentPageNo int    `json:"current_page_no"`
	TotalPages    int    `json:"total_pages"`
	PerPage       int    `json:"per_page"`
	TotalElements int64  `json:"total_elements"`
	Sort          string `json:"sort"`
}

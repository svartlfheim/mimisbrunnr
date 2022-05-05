package meta

type Pagination struct {
	Page  int `json:"page"`
	Limit int `json:"per_page"`
	Total int `json:"total"`
	Count int `json:"count"`
}

package entity

type Page struct {
	Page      int `json:"page"`
	PageSize  int `json:"pageSize"`
	PageCount int `json:"pageCount"`
}

func NewEmptyPage() Page {
	return Page{
		Page:      0,
		PageSize:  0,
		PageCount: 0,
	}
}

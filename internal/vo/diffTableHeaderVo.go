package vo

import "github.com/Catizard/lampghost_wails/internal/entity"

type DiffTableHeaderVo struct {
	DataUrl     string           `json:"data_url"`
	Name        string           `json:"name"`
	OriginalUrl *string          `json:"original_url"`
	Symbol      string           `json:"symbol"`
	Courses     [][]CourseInfoVo `json:"course"`
	HeaderUrl   string
}

func (header *DiffTableHeaderVo) Entity() *entity.DiffTableHeader {
	return &entity.DiffTableHeader{
		DataUrl:     header.DataUrl,
		Name:        header.Name,
		OriginalUrl: header.OriginalUrl,
		Symbol:      header.Symbol,
		HeaderUrl:   header.HeaderUrl,
	}
}

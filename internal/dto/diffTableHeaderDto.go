package dto

import (
	"github.com/Catizard/lampghost_wails/internal/entity"
)

type DiffTableHeaderDto struct {
  ID uint
	HeaderUrl   string
	DataUrl     string  
	Name        string  
	OriginalUrl *string 
	Symbol      string  

  Contents []entity.DiffTableData
}

func NewDiffTableHeaderDto(header *entity.DiffTableHeader, contents []entity.DiffTableData) *DiffTableHeaderDto {
  return &DiffTableHeaderDto{
    HeaderUrl: header.HeaderUrl,
    DataUrl: header.DataUrl,
    Name: header.Name,
    OriginalUrl: header.OriginalUrl,
    Symbol: header.Symbol,
    Contents: contents,
  }
}

package result

import "github.com/Catizard/lampghost_wails/internal/entity"

type RtnPage struct {
	Pagination entity.Page
	Rows       []interface{}
	RtnMessage
}

func NewRtnPage[T any](pagination entity.Page, rows []T) RtnPage {
	convRows := make([]interface{}, len(rows))
	for i := range convRows {
		convRows[i] = rows[i]
	}
	return RtnPage{
		pagination,
		convRows,
		SUCCESS,
	}
}

func NewErrorPage(err error) RtnPage {
	return RtnPage{
		entity.NewEmptyPage(),
		nil,
		NewErrorMessage(err),
	}
}

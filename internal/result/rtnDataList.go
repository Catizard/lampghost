package result

// NOTE: wails在v2之内还不能使用泛型
type RtnDataList struct {
	Rows []interface{}
	RtnMessage
}

func NewRtnDataList[T any](rows []T) RtnDataList {
	convRows := make([]interface{}, len(rows))
	for i := range convRows {
		convRows[i] = rows[i]
	}
	return RtnDataList{
		convRows,
		SUCCESS,
	}
}

func NewErrorDataList(err error) RtnDataList {
	return RtnDataList{
		nil,
		NewErrorMessage(err),
	}
}

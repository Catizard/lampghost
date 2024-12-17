package result

type RtnData struct {
	Data interface{}
	RtnMessage
}

func NewRtnData(data interface{}) RtnData {
	return RtnData{
		data,
		SUCCESS,
	}
}

func NewErrorData(err error) RtnData {
	return RtnData{
		nil,
		NewErrorMessage(err),
	}
}

package result

import "time"

type RtnMessage struct {
  Code int
  Msg string
  Timestamp time.Time
  Err error
}

var SUCCESS RtnMessage = NewRtnMessage(200, "Success")

func NewRtnMessage(code int, msg string) RtnMessage {
  return RtnMessage{
    Code: code,
    Msg: msg,
    Err: nil,
    Timestamp: time.Now(),
  }
}

func NewErrorMessage(err error) RtnMessage {
  return RtnMessage {
    Code: 500,
    Msg: err.Error(),
    Err: err,
    Timestamp: time.Now(),
  }
}

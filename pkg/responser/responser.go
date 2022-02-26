package responser

import "time"

var (
	GetLastTradeDataOk = add(1000, "已成功取得最新一筆資料")
)

func add(code int, msg string) ResponseFlag {
	return ResponseFlag{
		code: code, message: msg,
	}
}

func (responseFlag *ResponseFlag) Error() string {
	return responseFlag.message
}

func (responseFlag ResponseFlag) Message() string {
	return responseFlag.message
}

func (responseFlag ResponseFlag) Reload(message string) ResponseFlag {
	responseFlag.message = message
	return responseFlag
}

func (responseFlag ResponseFlag) Code() int {
	return responseFlag.code
}

type Response struct {
	ResultCode int         `json:"result_code"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
	TimeStamp  time.Time   `json:"time_stamp"`
}

type ResponseFlag struct {
	code    int
	message string
}

type ResponseFunc interface {
	Error() string
	Code() int
	Message() string
	Reload(string) ResponseFlag
}

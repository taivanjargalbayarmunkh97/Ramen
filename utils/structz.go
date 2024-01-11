package utils

type ResponseObj struct {
	ResponseCode int         `json:"response_code"`
	ResponseMsg  string      `json:"response_msg"`
	Data         interface{} `json:"data,omitempty"`
	Pagination   *Pagination `json:"pagination,omitempty"`
	Error        error       `json:"error,omitempty"`
}

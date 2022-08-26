package controller

const (
	errorNotAllowed         = "not_allowed"
	errorInvalidToken       = "invalid_token"
	errorSystemError        = "system_error"
	errorBadRequestBody     = "bad_request_body"
	errorBadRequestHeader   = "bad_request_header"
	errorBadRequestParam    = "bad_request_param"
	errorDuplicateCreating  = "duplicate_creating"
	errorResourceNotExists  = "resource_not_exists"
	errorConcurrentUpdating = "concurrent_updateing"
)

var (
	respBadRequestBody = newResponseCodeMsg(
		errorBadRequestBody, "can't fetch request body",
	)
)

// responseData is the response data to client
type responseData struct {
	Code string      `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func newResponseError(err error) responseData {
	code := errorSystemError

	return responseData{
		Code: code,
		Msg:  err.Error(),
	}
}

func newResponseData(data interface{}) responseData {
	return responseData{
		Data: data,
	}
}

func newResponseCodeError(code string, err error) responseData {
	return responseData{
		Code: code,
		Msg:  err.Error(),
	}
}

func newResponseCodeMsg(code, msg string) responseData {
	return responseData{
		Code: code,
		Msg:  msg,
	}
}

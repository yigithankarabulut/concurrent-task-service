package util

type ResponseData struct {
	Data   interface{} `json:"data"`
	Status int         `json:"status"`
}

type ErrorResponse struct {
	Error  string `json:"error"`
	Status int    `json:"status"`
}

func Response(Status int, Data interface{}) ResponseData {
	return ResponseData{
		Data:   Data,
		Status: Status,
	}
}

func BasicError(d interface{}, status int) ErrorResponse {
	var message string
	switch d.(type) {
	case error:
		message = d.(error).Error()
	case string:
		message = d.(string)
	case nil:
		message = "unknown error"
	default:
		message = "unknown error"
	}
	return ErrorResponse{
		Error:  message,
		Status: status,
	}
}

package Utils

type Response struct {
	Success bool
	Message string
	Body    interface{}
}

func FailedResponse(message string) Response {
	return Response{
		Success: false,
		Message: message,
		Body:    nil,
	}
}

func SucceededReponse(message string, body interface{}) Response {
	return Response{
		Success: true,
		Message: message,
		Body:    body,
	}
}

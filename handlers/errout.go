package handlers

type Response struct {
	Data  interface{}
	Good  bool
	Error string
}

func RespondData(data interface{}) Response {
	return Response{
		Data: data,
		Good: true,
	}
}

func ResponseError(err error) Response {
	if err == nil {
		return Response{
			Data:  nil,
			Good:  true,
			Error: "",
		}
	}
	return Response{
		Error: err.Error(),
		Good:  false,
	}
}

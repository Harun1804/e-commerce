package httpresponse

func ErrorDetail(err error) *interface{} {
	var details interface{} = err.Error()
	return &details
}
package assets_services

type ServiceError struct {
	Code int
	Err  error
}

func NewError(code int, err error) *ServiceError {
	return &ServiceError{Code: code, Err: err}
}
func (e *ServiceError) Error() string {
	return e.Err.Error()
}

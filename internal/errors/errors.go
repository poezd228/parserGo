package errors

type serviceErr struct {
	Field    string
	Err      error
	SiteCode string
}
type ServiceError interface {
	Error() error
}

func newErr(field string, err error) ServiceError {
	e := serviceErr{
		Field: field,
		Err:   err,
	}
	return &e
}

func WrongProxyUrl(err error) ServiceError {
	return newErr("Wrong proxy url", err)

}

func BadRequest(err error) ServiceError {
	return newErr("Bad request", err)

}
func UnableToCreateReq(err error) ServiceError {
	return newErr("UnableToCreateReq", err)

}
func UnableToUnmarshall(err error) ServiceError {
	return newErr("Unable To unmarshall Json ", err)
}
func (r *serviceErr) Error() error {
	return r.Err

}

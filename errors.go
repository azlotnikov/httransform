package httransform

import (
	"golang.org/x/xerrors"
)

var (
	// ErrProxyAuthorization is the error for ProxyAuthorizationBasicLayer
	// instance. If OnRequest callback of this method returns such an error,
	// then OnResponse callback generates correct 407 response.
	ErrProxyAuthorization = xerrors.New("cannot authenticate proxy user")
)

type RejectRequestError struct {
	responseText string
	responseCode int
}

func NewRejectRequestError(responseText string, responseCode int) *RejectRequestError {
	return &RejectRequestError{
		responseText: responseText,
		responseCode: responseCode,
	}
}

func (r *RejectRequestError) Error() string {
	return r.responseText
}

func IsRejectRequestError(err error) bool {
	_, ok := err.(*RejectRequestError)
	return ok
}

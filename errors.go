package httransform

import (
	"fmt"
	"golang.org/x/xerrors"
)

var (
	// ErrProxyAuthorization is the error for ProxyAuthorizationBasicLayer
	// instance. If OnRequest callback of this method returns such an error,
	// then OnResponse callback generates correct 407 response.
	ErrProxyAuthorization = xerrors.New("cannot authenticate proxy user")
)

type RejectRequestError struct {
	reason       string
	responseCode int
}

func NewRejectRequestError(reason string, responseCode int) *RejectRequestError {
	return &RejectRequestError{
		reason:       reason,
		responseCode: responseCode,
	}
}

func (r *RejectRequestError) Error() string {
	return fmt.Sprintf("request blocked with reason %s", r.reason)
}

func IsRejectRequestError(err error) bool {
	_, ok := err.(*RejectRequestError)
	return ok
}

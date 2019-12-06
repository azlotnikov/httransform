package ca

import "github.com/karlseguin/ccache"

type signRequest struct {
	host     string
	response chan signResponse
}

type signResponse struct {
	item ccache.TrackedItem
	err  error
}

package network

import "net/http"

func BuildTransport(skipVerifySSL bool) http.RoundTripper {
	return buildTransport(skipVerifySSL)
}

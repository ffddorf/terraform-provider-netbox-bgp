package client

import "net/http"

type HeadersTransport struct {
	http.RoundTripper
	headers map[string]string
}

func NewHeadersTransport(inner http.RoundTripper, headers map[string]string) *HeadersTransport {
	return &HeadersTransport{
		RoundTripper: inner,
		headers:      headers,
	}
}

func (ht *HeadersTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	for k, v := range ht.headers {
		req.Header.Set(k, v)
	}
	return ht.RoundTripper.RoundTrip(req)
}

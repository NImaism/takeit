package model

import (
	"crypto/tls"
	"fmt"
	"net/http"
)

type CancelTransport struct {
	Transport       http.RoundTripper
	TLSClientConfig *tls.Config
	MaxBodyBytes    int64
}

func (ct *CancelTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	resp, err := ct.Transport.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	if resp.ContentLength > ct.MaxBodyBytes {
		resp.Body.Close()
		return nil, fmt.Errorf("over")
	}

	return resp, nil
}

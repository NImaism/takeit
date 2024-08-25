package network

import (
	"crypto/tls"
	"github.com/nimaism/takeit/internal/model"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func InitHTTPClient(maxReadBody, timeoutSec int, verifySSL, disableRedirect bool) *http.Client {
	timeout := time.Duration(timeoutSec) * time.Second
	httpClient := http.Client{
		Timeout: timeout,
		Transport: &model.CancelTransport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: !verifySSL},
			Transport:       http.DefaultTransport,
			MaxBodyBytes:    int64(1000 * maxReadBody),
		},
	}

	if disableRedirect {
		httpClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}

	return &httpClient
}

func CheckValidURL(input string) bool {
	_, err := url.ParseRequestURI(input)
	return err == nil
}

func ExtractHost(urlString string) string {
	schemeEnd := strings.Index(urlString, "://")
	if schemeEnd == -1 {
		return ""
	}

	start := schemeEnd + 3

	end := strings.Index(urlString[start:], "/")
	if end == -1 {
		return urlString[start:]
	}

	return urlString[start : start+end]
}

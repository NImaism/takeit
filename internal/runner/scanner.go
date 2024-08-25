package runner

import (
	"fmt"
	"github.com/nimaism/takeit/pkg/network"
	"io"
	"net/http"
	"strings"
	"time"
)

func (r *Runner) Scan(target string) (bool, string, string, error) {
	host := target

	if !network.CheckValidURL(target) {
		target = "https://" + target
	} else {
		host = network.ExtractHost(target)
	}

	if r.Option.CheckCNAME {
		if cName, err := network.CheckCNAMERecord(host); err != nil {
			return false, "", "", fmt.Errorf("failed to check CNAME record: %v", err)
		} else if !cName {
			return false, "", "", nil
		}
	}

	if r.Option.Delay > 0 {
		time.Sleep(time.Duration(r.Option.Delay) * time.Second)
	}

	content, err := r.getContent(target)
	if err != nil {
		return false, "", "", fmt.Errorf("failed to get content: %v", err)
	}

	if found, discussion, documentation := r.matchContent(content); found {
		return true, documentation, discussion, nil
	}

	return false, "", "", nil
}

func (r *Runner) getContent(target string) (string, error) {
	var resp *http.Response
	var err error

	var attempts int

	req, reqErr := http.NewRequest("GET", target, nil)
	if reqErr != nil {
		return "	", fmt.Errorf("failed to create request: %v", reqErr)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:103.0) Gecko/20100101 Firefox/103.0")
	for k, v := range r.CustomHeader {
		req.Header.Set(k, v)
	}

	for attempts < r.Option.Retries {
		if !r.RateLimit.CanTake() {
			for !r.RateLimit.CanTake() {
				time.Sleep(time.Second)
			}
		}

		r.RateLimit.Take()

		resp, err = r.HttpClient.Do(req)
		if err == nil {
			break
		}

		// On failure, attempt with "http" if the initial request was made with "https"
		if req.URL.Scheme == "https" {
			req.URL.Scheme = "http"
			resp, err = r.HttpClient.Do(req)
			if err == nil {
				break
			}
		}

		attempts++
		if attempts < r.Option.Retries {
			time.Sleep(2 * time.Second)
		}
	}

	if err != nil {
		return "", fmt.Errorf("failed to get content after %d attempts: %v", attempts, err)
	}
	defer resp.Body.Close()

	content, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return "", fmt.Errorf("failed to read content: %v", readErr)
	}

	return string(content), nil
}

func (r *Runner) matchContent(content string) (bool, string, string) {
	for _, pattern := range *r.Patterns {
		if strings.Contains(content, pattern.Fingerprint) {
			found := true
			for _, falsePositiveString := range pattern.FalsePositive {
				if len(falsePositiveString) > 0 && strings.Contains(content, falsePositiveString) {
					found = false
					break
				}
			}
			if found {
				return true, pattern.Discussion, pattern.Documentation
			}
		}

	}

	return false, "", ""
}

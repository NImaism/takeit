package runner

import (
	"context"
	"fmt"
	"github.com/nimaism/takeit/internal/model"
	"github.com/nimaism/takeit/internal/options"
	"github.com/nimaism/takeit/internal/pattern"
	"github.com/nimaism/takeit/pkg/network"
	"github.com/projectdiscovery/ratelimit"
	"net/http"
	"sync"
	"time"
)

type Runner struct {
	Option       *options.Options
	Patterns     *[]model.Pattern
	RateLimit    *ratelimit.Limiter
	HttpClient   *http.Client
	CustomHeader map[string]string
}

func NewRunner(option *options.Options) (*Runner, error) {
	var limiter *ratelimit.Limiter

	if option.RateLimit > 0 {
		limiter = ratelimit.New(context.Background(), uint(option.RateLimit), time.Second)
	} else {
		limiter = ratelimit.New(context.Background(), uint(option.RateLimitMinute), time.Minute)
	}

	patterns, err := pattern.LoadPatterns(option.ExcludePatterns)
	if err != nil {
		return nil, fmt.Errorf("error loading patterns: %v", err)
	}

	return &Runner{
		Option:       option,
		HttpClient:   network.InitHTTPClient(option.BodyReadSize, option.Timeout, option.VerifySSL, option.DisableRedirects),
		Patterns:     patterns,
		RateLimit:    limiter,
		CustomHeader: option.ParseCustomHeaders(),
	}, nil
}

func (r *Runner) Run() {
	var reqChan = make(chan string)
	var wg sync.WaitGroup

	for i := 0; i < r.Option.Concurrency; i++ {
		wg.Add(1)
		go r.Worker(reqChan, &wg)
	}

	go func() {
		defer close(reqChan)
		for _, target := range r.Option.Targets {
			reqChan <- target
		}
	}()

	wg.Wait()
}

func (r *Runner) Worker(reqChan <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	for target := range reqChan {
		isVulnerable, discussion, documentation, err := r.Scan(target)

		if err != nil {
			r.printMessage("can't scan", target, false)
			continue
		}

		if isVulnerable {
			r.printScanResult(target, discussion, documentation)
		}
	}
}

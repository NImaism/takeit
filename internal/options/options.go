package options

import (
	"bufio"
	"fmt"
	"github.com/logrusorgru/aurora"
	"github.com/nimaism/takeit/internal/pattern"
	"github.com/projectdiscovery/goflags"
	"os"
	"strings"
)

type Options struct {
	// Targets is the list of targets to scan
	Targets goflags.StringSlice
	// VerifySSL verifies SSL certificates
	VerifySSL bool
	// DisableUpdateCheck disables automatic update check
	DisableUpdateCheck bool
	// CheckCNAME checks if the subdomain is a CNAME record
	CheckCNAME bool
	// Timeout is the timeout for the scan
	Timeout int
	// Concurrency is the number of concurrent goroutines
	Concurrency int
	// Silent shows only output
	Silent bool
	// CustomHeaders contains custom headers to send
	CustomHeaders goflags.StringSlice
	// BodyReadSize is the maximum size of response body to read
	BodyReadSize int
	// Delay is the delay between requests
	Delay int
	// Retries is the number of retries to do for network
	Retries int
	// ExcludePatterns are the patterns you don't want to check
	ExcludePatterns goflags.StringSlice
	//DisableRedirects disables the following of redirects
	DisableRedirects bool
	// RateLimit is the maximum number of requests to send per second
	RateLimit int
	// RateLimitMinute is the maximum number of requests to send per minute
	RateLimitMinute int
	// OutputFile is the file to write output to
	OutputFile string
	// NoColors disables coloring of response output
	NoColors bool
}

func (o *Options) ParseFlags() error {
	var cfgFile string

	flagSet := goflags.NewFlagSet()
	flagSet.SetDescription("Takeit is an advanced tool for detecting subdomain takeovers.")

	flagSet.CreateGroup("input", "Input",
		flagSet.StringSliceVarP(&o.Targets, "targets", "t", nil, "Targets to scan", goflags.FileCommaSeparatedStringSliceOptions),
	)

	flagSet.CreateGroup("config", "Configuration",
		flagSet.IntVarP(&o.BodyReadSize, "max-response-size", "mrs", 5*1000, "Maximum response size to read (kilobyte)"),
		flagSet.IntVar(&o.Timeout, "timeout", 10, "Time to wait for network in seconds"),
		flagSet.IntVar(&o.Retries, "retry", 1, "Number of times to retry the network"),
		flagSet.BoolVar(&o.VerifySSL, "verifySSL", false, "Verifies SSL certificates"),
		flagSet.StringVar(&cfgFile, "config", "", "Path to the configuration file"),
		flagSet.BoolVarP(&o.CheckCNAME, "cname", "cn", false, "Check CNAME before send request"),
		flagSet.StringSliceVarP(&o.CustomHeaders, "headers", "H", nil, "Custom header/cookie to include in all HTTP requests in header:value format (file)", goflags.FileStringSliceOptions),
		flagSet.StringSliceVarP(&o.ExcludePatterns, "exclude", "e", nil, "the patterns you don't want to check.", goflags.FileStringSliceOptions),
		flagSet.BoolVarP(&o.DisableRedirects, "disable-redirects", "dr", false, "Disable following redirects (default false)"),
	)

	flagSet.CreateGroup("ratelimit", "Rate-Limit",
		flagSet.IntVarP(&o.Concurrency, "concurrency", "c", 10, "Number of concurrent fetchers to use"),
		flagSet.IntVarP(&o.Delay, "delay", "rd", 0, "Request delay between each network in seconds"),
		flagSet.IntVarP(&o.RateLimit, "rate-limit", "rl", 150, "Maximum requests to send per second"),
		flagSet.IntVarP(&o.RateLimitMinute, "rate-limit-minute", "rlm", 0, "Maximum number of requests to send per minute"),
	)

	flagSet.CreateGroup("update", "Update",
		flagSet.BoolVarP(&o.DisableUpdateCheck, "disable-update-check", "duc", false, "Disable automatic update check"),
		flagSet.CallbackVarP(o.UpdatePatterns(), "update", "up", "update patterns to latest version"),
	)

	flagSet.CreateGroup("output", "Output",
		flagSet.BoolVarP(&o.NoColors, "no-color", "nc", false, "Disable output content coloring (ANSI escape codes)"),
		flagSet.BoolVar(&o.Silent, "silent", false, "Display output only"),
	)

	if err := flagSet.Parse(); err != nil {
		return fmt.Errorf("could not parse flags: %w", err)
	}

	if cfgFile != "" {
		if err := flagSet.MergeConfigFile(cfgFile); err != nil {
			return fmt.Errorf("could not parse configuration file: %w", err)
		}
	}

	if len(o.Targets) < 1 {
		fi, err := os.Stdin.Stat()
		if err != nil {
			return fmt.Errorf("could not read from stdin")
		}

		if fi.Mode()&os.ModeNamedPipe == 0 {
			return nil
		}

		sc := bufio.NewScanner(os.Stdin)
		for sc.Scan() {
			if sc.Text() != "" {
				o.Targets = append(o.Targets, sc.Text())
			}
		}

		if err = sc.Err(); err != nil {
			return fmt.Errorf("could not parse stdin")
		}
	}

	return nil
}

func (o *Options) ParseCustomHeaders() map[string]string {
	customHeaders := make(map[string]string)
	for _, header := range o.CustomHeaders {
		if parts := strings.SplitN(header, ":", 2); len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			customHeaders[key] = value
		}
	}
	return customHeaders
}

func (o *Options) UpdatePatterns() func() {
	return func() {
		if err := pattern.UpdatePatterns(o.ExcludePatterns); err != nil {
			fmt.Println(aurora.Red(err.Error()))
		}
		fmt.Println(aurora.Green("Patterns updated successfully"))
	}
}

package worker

import (
	"context"
	"net/http"
	"time"

	"github.com/pulkitbhatt/ikiru/internal/queue"
)

type ResultStatus string

const (
	MaxRetryCount = 2

	StatusSuccess ResultStatus = "success"
	StatusFailure ResultStatus = "failure"

	StatusCodeEmpty = 0

	ErrConnectionRefused = "connection_refused"
	ErrTimeout           = "timeout"
	ErrConnection        = "connection_error"
	ErrDNS               = "dns_error"
	ErrNetwork           = "network_error"
	ErrInvalidRequest    = "invalid_request"
	ErrNon2XXStatus      = "non_2xx_status_code"
)

type CheckResult struct {
	Status     ResultStatus
	HTTPStatus int
	LatencyMs  int
	Error      string
}

var defaultHTTPClient = &http.Client{
	Transport: &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     30 * time.Second,
	},
}

func ExecuteHTTPCheck(ctx context.Context, job queue.MonitorJob) CheckResult {
	start := time.Now()

	ctx, cancel := context.WithTimeout(ctx, time.Duration(job.TimeoutMs)*time.Millisecond)
	defer cancel()

	var lastErr error

	for attempt := 0; attempt <= MaxRetryCount; attempt++ {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, job.URL, nil)
		if err != nil {
			return failResult(ErrInvalidRequest,
				err,
				time.Since(start).Milliseconds(),
				StatusCodeEmpty,
			)
		}

		resp, err := defaultHTTPClient.Do(req)
		latency := time.Since(start).Milliseconds()

		if err != nil {
			lastErr = err

			if ctx.Err() == context.DeadlineExceeded {
				return failResult(ErrTimeout, ctx.Err(), latency, StatusCodeEmpty)
			}

			time.Sleep(time.Duration(attempt+1) * 100 * time.Millisecond)
			continue
		}

		defer resp.Body.Close()

		if resp.StatusCode >= 200 && resp.StatusCode < 400 {
			return successResult(latency, resp.StatusCode)
		}

		return CheckResult{
			Status:     StatusFailure,
			LatencyMs:  int(latency),
			HTTPStatus: resp.StatusCode,
			Error:      ErrNon2XXStatus,
		}
	}
	return CheckResult{Error: lastErr.Error()}
}

func successResult(latency int64, statusCode int) CheckResult {
	return CheckResult{
		Status:     StatusSuccess,
		LatencyMs:  int(latency),
		HTTPStatus: statusCode,
	}
}

func failResult(kind string, err error, latency int64, statusCode int) CheckResult {
	return CheckResult{
		Status:     StatusFailure,
		LatencyMs:  int(latency),
		Error:      kind + ": " + err.Error(),
		HTTPStatus: statusCode,
	}
}

package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/ryanfaerman/misty/config"
	log "github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
)

const (
	userAgent = "go-mistyclient"
)

var (
	Version = "1"

	DefaultDialer                      = &net.Dialer{Timeout: 1000 * time.Millisecond}
	DefaultTransport http.RoundTripper = &http.Transport{Dial: DefaultDialer.Dial, Proxy: http.ProxyFromEnvironment}
	DefaultClient                      = &http.Client{Transport: DefaultTransport}

	DefaultRequestTimeout = 5 * time.Second
	DefaultRetries        = 3
	DefaultRetryDelay     = 3 * time.Second
	DefaultRequestRate    = 100
	DefaultRequestBurst   = 15
)

func SetConnectTimeout(duration time.Duration) {
	DefaultDialer.Timeout = duration
}

type service struct {
	client *Client
}

// Client is the thing in charge of interacting with SPADE. It's basically
// an HTTP client with some extra bells and whistles. It supports rate limiting.
type Client struct {
	BaseURL   *url.URL
	UserAgent string
	Username  string
	Password  string

	Retries        int
	RetryDelay     time.Duration
	RequestRate    int
	RequestBurst   int
	RequestTimeout time.Duration

	limiter  *rate.Limiter
	client   *http.Client
	clientMu sync.Mutex

	logger *log.Logger

	common service // reuse a single struct instead of allocating one fore ach service on the heap.

	// Services used for talking to different parts of the KB API
	Display *DisplayService
}

// New returns a new Client ready for use. If httpClient is nil,
// it will just use http.DefaultClient.
func New(addr string, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = DefaultClient
	}

	if !strings.HasSuffix(addr, "/") {
		addr = addr + "/"
	}
	baseURL, _ := url.Parse(addr)

	c := &Client{
		client:         httpClient,
		BaseURL:        baseURL,
		UserAgent:      userAgent,
		Retries:        DefaultRetries,
		RetryDelay:     DefaultRetryDelay,
		RequestRate:    DefaultRequestRate,
		RequestBurst:   DefaultRequestBurst,
		RequestTimeout: DefaultRequestTimeout,

		logger: config.Logger,
	}
	c.common.client = c

	c.Display = (*DisplayService)(&c.common)
	// c.Page = (*PageService)(&c.common)
	// c.Search = (*SearchService)(&c.common)

	return c
}

// Do sends an HTTP request and returns an HTTP response, following the
// settings configured by the client.
func (client *Client) Do(ctx context.Context, req *http.Request) (*http.Response, error) {
	if client.limiter == nil {
		client.limiter = rate.NewLimiter(rate.Limit(client.RequestRate), client.RequestBurst)
	}

	if _, ok := ctx.Deadline(); !ok {
		ctxWithTimeout, cancel := context.WithTimeout(ctx, 4*(client.RetryDelay+client.RequestTimeout))
		ctx = ctxWithTimeout
		defer cancel()
	}

	if err := client.limiter.Wait(ctx); err != nil {
		return nil, err
	}

	if _, _, ok := req.BasicAuth(); !ok {
		req.SetBasicAuth(client.Username, client.Password)
	}

	req.Header.Set("User-Agent", fmt.Sprintf("%s/%s", userAgent, Version))
	req.Header.Set("Accept", "application/json")

	// if req.Body != nil {
	req.Header.Set("Content-Type", "application/json")
	// }

	var res *http.Response
	var err error

	res, err = client.client.Do(req)

	for i := 1; err != nil && i < client.Retries; i++ {
		select {
		case <-time.After(client.RetryDelay):
			res, err = client.client.Do(req)
			continue
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}

	return res, err
}

func (client *Client) NewRequest(method string, urlStr string, body interface{}) (*http.Request, error) {
	if !strings.HasSuffix(client.BaseURL.Path, "/") {
		return nil, fmt.Errorf("BaseURL must have a trailing slash, but %q does not", client.BaseURL)
	}
	u, err := client.BaseURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		err := enc.Encode(body)
		if err != nil {
			return nil, err
		}
	}

	return http.NewRequest(method, u.String(), buf)
}

func (client *Client) Get(ctx context.Context, url string) (*http.Response, error) {
	req, err := client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	return client.Do(ctx, req)
}

func (client *Client) Post(ctx context.Context, url string, body interface{}) (*http.Response, error) {
	req, err := client.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}

	return client.Do(ctx, req)
}

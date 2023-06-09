package surf

import (
	"context"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/x0xO/hg/surf/pkg/drainbody"
)

// Request a struct that holds information about an HTTP request.
type Request struct {
	request    *http.Request
	client     *Client
	writeErr   *error
	error      error
	remoteAddr net.Addr
	body       io.ReadCloser
}

func (req *Request) GetRequest() *http.Request { return req.request }

// Do performs the HTTP request and returns a Response object or an error if the request failed.
func (req *Request) Do() (*Response, error) {
	if req.error != nil {
		return nil, req.error
	}

	if err := req.client.applyReqMW(req); err != nil {
		return nil, err
	}

	opt := req.client.opt
	if opt != nil {
		if err := opt.applyReqMW(req); err != nil {
			return nil, err
		}
	}

	// clone request body
	req.body, req.request.Body, req.error = drainbody.DrainBody(req.request.Body)
	if req.error != nil {
		return nil, req.error
	}

	var (
		resp     *http.Response
		attempts int
		err      error
	)

	start := time.Now()
	cli := req.client.cli

	for {
		resp, err = cli.Do(req.request)

		notRetriable := err == nil &&
			resp.StatusCode != http.StatusInternalServerError &&
			resp.StatusCode != http.StatusTooManyRequests &&
			resp.StatusCode != http.StatusServiceUnavailable

		if notRetriable || opt == nil || opt.retryMax == 0 || attempts >= opt.retryMax {
			break
		}

		attempts++

		time.Sleep(opt.retryWait)
	}

	if err != nil {
		return nil, err
	}

	if req.writeErr != nil && (*req.writeErr).Error() != "" {
		return nil, *req.writeErr
	}

	response := &Response{
		Attempts:      attempts,
		Time:          time.Since(start),
		Client:        req.client,
		ContentLength: resp.ContentLength,
		Cookies:       resp.Cookies(),
		Headers:       headers(resp.Header),
		History:       req.client.history,
		Proto:         resp.Proto,
		Status:        resp.Status,
		StatusCode:    resp.StatusCode,
		URL:           resp.Request.URL,
		UserAgent:     req.request.UserAgent(),
		remoteAddr:    req.remoteAddr,
		request:       req,
		response:      resp,
		Body: &body{
			body:        resp.Body,
			cache:       opt != nil && opt.cacheBody,
			contentType: resp.Header.Get("Content-Type"),
			deflate:     resp.Header.Get("Content-Encoding") == "deflate",
			limit:       -1,
		},
	}

	if err := req.client.applyRespMW(response); err != nil {
		return nil, err
	}

	return response, nil
}

// WithContext associates the provided context with the request.
func (req *Request) WithContext(ctx context.Context) *Request {
	if ctx != nil {
		req.request = req.request.WithContext(ctx)
	}

	return req
}

// AddCookies adds cookies to the request.
func (req *Request) AddCookies(cookies ...http.Cookie) *Request {
	for _, cookie := range cookies {
		c := cookie
		req.request.AddCookie(&c)
	}

	return req
}

// SetHeaders sets headers for the request, replacing existing ones with the same name.
func (req *Request) SetHeaders(headers map[string]string) *Request {
	if headers != nil && req.request != nil {
		for header, data := range headers {
			req.request.Header.Set(header, data)
		}
	}

	return req
}

// AddHeaders adds headers to the request, appending to any existing headers with the same name.
func (req *Request) AddHeaders(headers map[string]string) *Request {
	if headers != nil && req.request != nil {
		for header, data := range headers {
			req.request.Header.Add(header, data)
		}
	}

	return req
}

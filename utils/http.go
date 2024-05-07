package utils

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type Request struct {
	url        string
	httpClient *http.Client
	headers    map[string]string
	cookies    map[string]string
	query      map[string]string
	body       []byte
}

type Response struct {
	res  *http.Response
	body []byte
}

func NewRequest(url string) *Request {
	httpClient := &http.Client{}
	request := &Request{url: url, httpClient: httpClient}
	return request
}

func (r *Request) Headers(headers map[string]string) *Request {
	r.headers = headers
	return r
}

func (r *Request) Cookies(cookies map[string]string) *Request {
	r.cookies = cookies
	return r
}

func (r *Request) Body(body []byte) *Request {
	r.body = body
	return r
}

func (r *Response) Body() []byte {
	return r.body
}

func (r *Response) Status() int {
	return r.res.StatusCode
}

func (r *Response) Unmarshal(v any) error {
	return json.Unmarshal(r.body, &v)
}

func (r Request) Get(endpoint string) (*Response, error) {
	req, err := http.NewRequest(http.MethodGet, r.url+endpoint, nil)
	if err != nil {
		return nil, err
	}
	prepReq := r.prepareRequest(req)
	return r.sendRequest(prepReq)
}

func (r Request) Post(endpoint string) (*Response, error) {
	req, err := http.NewRequest(http.MethodPost, r.url+endpoint, bytes.NewBuffer(r.body))
	if err != nil {
		return nil, err
	}
	prepReq := r.prepareRequest(req)
	return r.sendRequest(prepReq)
}

func (r Request) Put(endpoint string) (*Response, error) {
	req, err := http.NewRequest(http.MethodPut, r.url+endpoint, bytes.NewBuffer(r.body))
	if err != nil {
		return nil, err
	}
	prepReq := r.prepareRequest(req)
	return r.sendRequest(prepReq)
}

func (r *Request) prepareRequest(req *http.Request) *http.Request {
	// Set Headers
	for key, value := range r.headers {
		req.Header.Set(key, value)
	}
	// Set Queries
	q := req.URL.Query()
	for key, value := range r.query {
		q.Add(key, value)
	}
	req.URL.RawQuery = q.Encode()
	return req
}

func (r *Request) sendRequest(req *http.Request) (*Response, error) {
	res, err := r.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return &Response{res, body}, nil
}

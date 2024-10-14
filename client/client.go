package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	host       string
	headers    map[string]string
	http       http.Client
	intercepts []func(*Client)
	handlers   map[int]func()
}

type ClientInterface interface {
	Get(endpoint string, rv any) error
	Post(endpoint string, data any, rv any) error
}

// Do a GET request to an endpoint. rv is used to unmarshal the result to any given GO value
func (c *Client) Get(endpoint string, rv any) error {
	for _, i := range c.intercepts {
		i(c)
	}

	req, readErr := http.NewRequest("GET", c.host+endpoint, nil)
	if readErr != nil {
		return readErr
	}

	for h := range c.headers {
		req.Header.Add(h, c.headers[h])
	}

	res, resErr := c.http.Do(req)
	if resErr != nil {
		return resErr
	}

	resBody, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		return readErr
	}

	if h, ok := c.handlers[res.StatusCode]; ok {
		h()
	}

	if res.StatusCode >= 300 {
		return fmt.Errorf("%s :: %+v", res.Status, string(resBody[:]))
	}

	if rv != nil {
		jsonErr := json.Unmarshal(resBody, &rv)
		if jsonErr != nil {
			return jsonErr
		}
	}
	return nil
}

// Do a POST request to an endpoint. rv is used to unmarshal the result to any given GO value
func (c *Client) Post(endpoint string, data any, rv any) error {
	bodyBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	for _, i := range c.intercepts {
		i(c)
	}

	req, reqErr := http.NewRequest("POST", c.host+endpoint, bytes.NewReader(bodyBytes))
	if reqErr != nil {
		return reqErr
	}

	for h := range c.headers {
		req.Header.Add(h, c.headers[h])
	}

	res, resErr := c.http.Do(req)
	if resErr != nil {
		return resErr
	}

	resBody, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		return readErr
	}

	if h, ok := c.handlers[res.StatusCode]; ok {
		h()
	}

	if res.StatusCode >= 300 {
		return fmt.Errorf("%s :: %+v", res.Status, string(resBody[:]))
	}

	if rv != nil {
		jsonErr := json.Unmarshal(resBody, &rv)
		if jsonErr != nil {
			return jsonErr
		}
	}
	return nil
}

// Initialize a new client
func NewClient(host string, options ...func(*Client)) *Client {
	client := &Client{
		host:     host,
		headers:  make(map[string]string),
		http:     http.Client{},
		handlers: make(map[int]func()),
	}
	for _, o := range options {
		o(client)
	}
	return client
}

// Add options outside of initialization
func (c *Client) Config(options ...func(*Client)) {
	for _, o := range options {
		o(c)
	}
}

// Explicitly set a header on the client
func (c *Client) SetHeader(key string, value string) {
	c.headers[key] = value
}

// Add a header to the client that will be applied before every request
func WithHeader(key string, value string) func(*Client) {
	return func(c *Client) {
		c.headers[key] = value
	}
}

// Shortcut for WithHeader("Authorization", <token>)
func WithAuth(value string) func(*Client) {
	return func(c *Client) {
		c.headers["Authorization"] = value
	}
}

// Add a function that will run before every method
//
// Useful when dealing with headers that cannot be set on initialization
func WithInterceptor(i func(*Client)) func(*Client) {
	return func(c *Client) {
		c.intercepts = append(c.intercepts, i)
	}
}

// Add a function that will run if the status code in any response matches statusCode
func WithStatusHandler(statusCode int, h func()) func(*Client) {
	return func(c *Client) {
		c.handlers[statusCode] = h
	}
}

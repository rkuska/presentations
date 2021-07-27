package client

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"
)

type Client struct {
	c *http.Client
}

func NewClient() *Client {
	c := http.Client{
		Timeout: 60 * time.Second,
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 15 * time.Second,
			}).Dial,
			ResponseHeaderTimeout: 10 * time.Second,
			IdleConnTimeout:       6 * time.Minute,
			MaxIdleConns:          100,
			MaxIdleConnsPerHost:   100,
		},
	}
	return &Client{c: &c}
}

func (c *Client) GetSleepy(ctx context.Context) {
	req, err := http.NewRequestWithContext(ctx, "GET", "http://127.0.0.1:8080/sleepyget", nil)
	if err != nil {
		fmt.Println("failed to create request")
		return
	}
	resp, err := c.c.Do(req)
	if err != nil {
		fmt.Println("received error from request:", err.Error())
		return
	}
	defer resp.Body.Close()
	fmt.Println("received response: ", resp.StatusCode, resp.Proto)
}

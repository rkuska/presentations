package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

const (
	// extraDataMaxBytes maximum number of bytes to read after `More()`.
	extraDataMaxBytes = 4096
)

type Client struct {
	c *http.Client
}

type response struct {
	Msg string
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
	var r response
	err = unmarshalReader(resp.Body, &r)
	if err != nil {
		fmt.Println("failed to unmarshal the response: %s", err.Error())
		return
	}
	fmt.Println("received response: ", resp.StatusCode, resp.Proto)
}

// unmarshalReader is a helper method to decode the body of the response. Once it manages to decode
// the body it tries to call `More` to make sure there are no additional data in the response. If
// yes, it tries to read them and returns them as error containing the additional data as a string.
func unmarshalReader(r io.Reader, val interface{}) error {
	decoder := json.NewDecoder(r)
	err := decoder.Decode(val)
	if err != nil {
		return err
	}
	if decoder.More() {
		extraDataReader := io.LimitReader(io.MultiReader(decoder.Buffered(), r), extraDataMaxBytes)
		extraData, err := ioutil.ReadAll(extraDataReader)
		if err != nil {
			return err
		}
		return errors.New(fmt.Sprintf("there are more data after the response was read: %s", extraData))
	}
	return nil
}

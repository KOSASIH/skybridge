package network

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
 "time"

	"github.com/KOSASIH/skybridge/utils"
)

// Network represents a network connection
type Network struct {
	httpClient *http.Client
}

// NewNetwork returns a new Network instance
func NewNetwork() *Network {
	return &Network{
		httpClient: &http.Client{
			Transport: &http.Transport{
				DialContext: (&net.Dialer{
					Timeout: 30 * time.Second,
				}).DialContext,
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		},
	}
}

// Get performs a GET request to the specified URL
func (n *Network) Get(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := n.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code %d", resp.StatusCode)
	}

	body, err := utils.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

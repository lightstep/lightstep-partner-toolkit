package backstageprocessor

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"go.uber.org/zap"
)

// Client defines the Backstage client interface
type Client interface {
	GetEntities() (*BackstageEntityResponse, error)
}

// NewClientProvider creates the default rest client provider
func NewClientProvider(endpoint url.URL, logger *zap.Logger) ClientProvider {
	return &defaultClientProvider{
		endpoint: endpoint,
		logger:   logger,
	}
}

// ClientProvider defines
type ClientProvider interface {
	BuildClient() Client
}

type defaultClientProvider struct {
	endpoint url.URL
	logger   *zap.Logger
}

func (dcp *defaultClientProvider) BuildClient() Client {
	return defaultClient(
		dcp.endpoint,
		dcp.logger,
	)
}

// TODO: Try using config.HTTPClientSettings
func defaultClient(
	endpoint url.URL,
	logger *zap.Logger,
) *clientImpl {
	tr := defaultTransport()
	return &clientImpl{
		baseURL:    endpoint,
		httpClient: http.Client{Transport: tr},
		logger:     logger,
	}
}

func defaultTransport() *http.Transport {
	return http.DefaultTransport.(*http.Transport).Clone()
}

var _ Client = (*clientImpl)(nil)

type clientImpl struct {
	baseURL    url.URL
	httpClient http.Client
	logger     *zap.Logger
}

func (c *clientImpl) get(path string) ([]byte, error) {
	req, err := c.buildReq(path)
	if err != nil {
		return nil, err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		closeErr := resp.Body.Close()
		if closeErr != nil {
			c.logger.Warn("Failed to close response body", zap.Error(closeErr))
		}
	}()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request GET %s failed - %q", req.URL.String(), resp.Status)
	}
	return body, nil
}

func (c *clientImpl) GetEntities() (*BackstageEntityResponse, error) {
	resp, err := c.get("/api/catalog/entities?filter=kind=Component")
	if err != nil {
		return nil, err
	}

	var backstageResp BackstageEntityResponse
	json.Unmarshal(resp, &backstageResp)
	return &backstageResp, nil
}

func (c *clientImpl) buildReq(path string) (*http.Request, error) {
	url := c.baseURL.String() + path
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

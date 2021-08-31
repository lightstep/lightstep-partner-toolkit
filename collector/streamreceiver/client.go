package streamreceiver

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"go.uber.org/zap"
)

// Client defines the Lightstep API client interface
type Client interface {
	GetStreamTraces() (*LightstepStreamResponse, error)
}

// NewClientProvider creates the default rest client provider
func NewClientProvider(endpoint url.URL, org string, project string, apiKey string, logger *zap.Logger) ClientProvider {
	return &defaultClientProvider{
		endpoint: endpoint,
		logger:   logger,
		org: org,
		project: project,
		apiKey: apiKey,
	}
}

// ClientProvider defines
type ClientProvider interface {
	BuildClient() Client
}

type defaultClientProvider struct {
	endpoint url.URL
	org string
	project string
	apiKey string
	logger   *zap.Logger
}

func (dcp *defaultClientProvider) BuildClient() Client {
	return defaultClient(
		dcp.endpoint,
		dcp.logger,
		dcp.org,
		dcp.project,
		dcp.apiKey,
	)
}

// TODO: Try using config.HTTPClientSettings
func defaultClient(
	endpoint url.URL,
	logger *zap.Logger,
	org string,
	project string,
	apiKey string,
) *clientImpl {
	tr := defaultTransport()
	return &clientImpl{
		baseURL:    endpoint,
		httpClient: http.Client{Transport: tr},
		org: org,
		project: project,
		apiKey: apiKey,
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
	apiKey string
	project string
	org string
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

func (c *clientImpl) GetStreamTraces() (*LightstepStreamResponse, error) {
	path := fmt.Sprintf("%s/projects/%s/streams/QQRrtJmW/timeseries?resolution-ms=60000&youngest-time=2021-08-30T17:58:07.177Z&oldest-time=2021-08-30T17:00:07.177Z&include-exemplars=1", c.org, c.project)
	resp, err := c.get(path)
	if err != nil {
		return nil, err
	}

	var lsResp LightstepStreamResponse
	json.Unmarshal(resp, &lsResp)
	return &lsResp, nil
}

func (c *clientImpl) buildReq(path string) (*http.Request, error) {
	url := c.baseURL.String() + path
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("bearer %s", c.apiKey))
	return req, nil
}
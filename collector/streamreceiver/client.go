package streamreceiver

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"go.uber.org/zap"
)

// Client defines the Lightstep API client interface
type Client interface {
	GetStreamTraces() (*LightstepStreamResponse, error)
	GetTrace(string) (*LightstepTraceResponse, error)

}

// NewClientProvider creates the default rest client provider
func NewClientProvider(endpoint url.URL, org string, project string, apiKey string, streamId string, logger *zap.Logger) ClientProvider {
	return &defaultClientProvider{
		endpoint: endpoint,
		logger:   logger,
		org: org,
		project: project,
		apiKey: apiKey,
		streamId: streamId,
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
	streamId string
	logger   *zap.Logger
}

func (dcp *defaultClientProvider) BuildClient() Client {
	return defaultClient(
		dcp.endpoint,
		dcp.logger,
		dcp.org,
		dcp.project,
		dcp.apiKey,
		dcp.streamId,
	)
}

// TODO: Try using config.HTTPClientSettings
func defaultClient(
	endpoint url.URL,
	logger *zap.Logger,
	org string,
	project string,
	apiKey string,
	streamId string,
) *clientImpl {
	tr := defaultTransport()
	return &clientImpl{
		baseURL:    endpoint,
		httpClient: http.Client{Transport: tr},
		org: org,
		project: project,
		apiKey: apiKey,
		streamId: streamId,
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
	streamId string
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

func (c *clientImpl) GetTrace(spanId string) (*LightstepTraceResponse, error) {
	path := fmt.Sprintf("%s/projects/%s/stored-traces?span-id=%s", c.org, c.project, spanId)
	resp, err := c.get(path)
	if err != nil {
		return nil, err
	}
	var lsResp LightstepTraceResponse
	err = json.Unmarshal(resp, &lsResp); if err != nil {
		return nil, err
	}

	return &lsResp, nil
}

func (c *clientImpl) GetStreamTraces() (*LightstepStreamResponse, error) {
	youngestTime := time.Now().Add(time.Duration(-5) * time.Minute)
	oldestTime := youngestTime.Add(time.Duration(-3) * time.Minute)

	path := fmt.Sprintf("%s/projects/%s/streams/%s/timeseries?resolution-ms=60000&youngest-time=%s&oldest-time=%s&include-exemplars=1", c.org, c.project, c.streamId, youngestTime.Format(time.RFC3339), oldestTime.Format(time.RFC3339))
	resp, err := c.get(path)
	if err != nil {
		return nil, err
	}

	var lsResp LightstepStreamResponse
	err = json.Unmarshal(resp, &lsResp); if err != nil {
		return nil, err
	}
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
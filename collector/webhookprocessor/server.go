// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package webhookprocessor

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/model/pdata"
	"go.uber.org/zap"
)

type httpServer struct {
	serverType string
	server     *http.Server
	logger     *zap.Logger
	config     *Config
	attrProc   *AttrProc
	actions    []ActionKeyValue
}

const (
	TracesServer  = "traces"
	MetricsServer = "metrics"
)

func (h *httpServer) Start(_ context.Context, host component.Host) error {

	handler := http.NewServeMux()
	handler.HandleFunc("/webhook", h.webhookHandler())
	handler.HandleFunc("/upsert", h.actionHandler(UPSERT))
	handler.HandleFunc("/delete", h.actionHandler(DELETE))

	var listener net.Listener
	var err error
	if h.serverType == MetricsServer {
		h.server = h.config.MetricsIngress.ToServer(handler)
		listener, err = h.config.MetricsIngress.ToListener()
		if err != nil {
			return fmt.Errorf("failed to bind to address %s: %w", h.config.MetricsIngress.Endpoint, err)
		}
	} else if h.serverType == TracesServer {
		h.server = h.config.TracesIngress.ToServer(handler)
		listener, err = h.config.TracesIngress.ToListener()
	} else {
		host.ReportFatalError(fmt.Errorf("could not identify server type"))
	}

	if err != nil {
		return fmt.Errorf("failed to bind to address %s: %w", h.config.TracesIngress.Endpoint, err)
	}
	go func() {
		if err := h.server.Serve(listener); err != http.ErrServerClosed {
			host.ReportFatalError(err)
		}
	}()

	err = h.setAttrProc()
	if err != nil {
		return err
	}

	return nil
}
func (h *httpServer) addAttrAction(actionKeyValue ActionKeyValue) error {
	h.actions = append(h.actions, actionKeyValue)
	err := h.setAttrProc()
	if err != nil {
		return fmt.Errorf("could not add ActionKeyValue: %v", err)
	}
	return nil
}

func (h *httpServer) setAttrProc() error {
	attrProc := &AttrProc{
		actions: h.actions,
	}

	h.attrProc = attrProc
	return nil
}

func (h *httpServer) ProcessMetrics(_ context.Context, md pdata.Metrics) (pdata.Metrics, error) {
	rl := md.ResourceMetrics()
	for i := 0; i < rl.Len(); i++ {
		resource := rl.At(i).Resource()
		h.attrProc.Process(resource.Attributes(), "")
	}
	return md, nil
}

func (h *httpServer) ProcessTraces(_ context.Context, td pdata.Traces) (pdata.Traces, error) {
	rss := td.ResourceSpans()
	for i := 0; i < rss.Len(); i++ {
		rs := rss.At(i)
		serviceNameValue, _ := rs.Resource().Attributes().Get("service.name")
		serviceName := serviceNameValue.StringVal()

		ilss := rs.InstrumentationLibrarySpans()
		for j := 0; j < ilss.Len(); j++ {
			ils := ilss.At(j)
			spans := ils.Spans()
			for k := 0; k < spans.Len(); k++ {
				span := spans.At(k)

				h.attrProc.Process(span.Attributes(), serviceName)
			}
		}
	}
	return td, nil
}

func (h *httpServer) Shutdown(ctx context.Context) error {
	if err := h.server.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}

func (h *httpServer) removeAttribute(key string, service string) error {
	actionKeyValue := ActionKeyValue{
		Key:    key,
		Action: DELETE,
		ServiceName: service,
	}
	h.logger.Debug("removing attribute", zap.String("key", key))
	return h.addAttrAction(actionKeyValue)
}

func (h *httpServer) addAttribute(key string, value string, service string) error {
	actionKeyValue := ActionKeyValue{
		Key:    key,
		Value:  value,
		Action: UPSERT,
		ServiceName: service,
	}
	h.logger.Debug("adding attribute",  zap.String("key", key))
	return h.addAttrAction(actionKeyValue)
}

func (h *httpServer) actionHandler(actionType Action) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		k := r.URL.Query().Get("key")
		v := r.URL.Query().Get("value")
		service := r.URL.Query().Get("service")
		//from := r.URL.Query().Get("from_attribute")
		if len(k) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "bad request: expected key param")
			return
		}

		var err error
		if actionType == DELETE {
			err = h.removeAttribute(k, service)
		} else if actionType == UPSERT {
			err = h.addAttribute(k, v, service)
		}

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "internal error: could not set attr proc: %v", err)
			return
		}
		w.WriteHeader(http.StatusAccepted)
		h.logger.Info("added new processor rule", zap.String("key", k), zap.String("value", v), zap.String("service", service))
		fmt.Fprintf(w, "ok: %v", actionType)
	}
}

func newHTTPServer(config *Config, logger *zap.Logger, serverType string) (*httpServer, error) {
	h := &httpServer{
		config:     config,
		logger:     logger,
		serverType: serverType,
	}

	h.actions = []ActionKeyValue{}

	return h, nil
}

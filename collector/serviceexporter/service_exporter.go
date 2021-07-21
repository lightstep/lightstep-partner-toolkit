package serviceexporter

import (
	"context"
	"encoding/json"
	"fmt"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/consumer/pdata"
	"log"
	"net/http"
	"sync"
)

type ServiceResourceAttributes map[string]interface{}

type ServiceResources struct  {
	Services map[string]ServiceResourceAttributes `json:"services"`
}

type serviceExporter struct {
	mutex sync.Mutex
	server     *http.Server
	config     *Config
	serviceResources *ServiceResources
}

func (e *serviceExporter) Capabilities() consumer.Capabilities {
	return consumer.Capabilities{MutatesData: false}
}

func (e *serviceExporter) ConsumeTraces(_ context.Context, td pdata.Traces) error {
	rss := td.ResourceSpans()
	for i := 0; i < rss.Len(); i++ {
		rs := rss.At(i)
		serviceName, ok := rs.Resource().Attributes().Get("service.name")
		if ok {
			resourceAttrs := attrsValue(rs.Resource().Attributes())
			e.serviceResources.Services[serviceName.StringVal()] = resourceAttrs
			log.Printf("resourceAttrs for service %v: %v", serviceName, resourceAttrs)
		}
	}
	return nil
}

func (e *serviceExporter) Start(_ context.Context, host component.Host) error {
	handler := http.NewServeMux()

	handler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		data, _ := json.Marshal(e.serviceResources)
		_, _ = fmt.Fprintf(w, "%s", string(data))
	})

	e.server = e.config.Scraper.ToServer(handler)
	listener, err := e.config.Scraper.ToListener()
	if err != nil {
		return fmt.Errorf("failed to bind to address %s: %w", e.config.Scraper.Endpoint, err)
	}

	go func() {
		if err := e.server.Serve(listener); err != nil {
			host.ReportFatalError(err)
		}
	}()
	return nil
}

func (e *serviceExporter) Shutdown(context.Context) error {
	return e.server.Close()
}

func attrsValue(attrs pdata.AttributeMap) map[string]interface{} {
	if attrs.Len() == 0 {
		return nil
	}
	out := make(map[string]interface{}, attrs.Len())
	attrs.Range(func(k string, v pdata.AttributeValue) bool {
		out[k] = attrValue(v)
		return true
	})
	return out
}

func attrValue(value pdata.AttributeValue) interface{} {
	switch value.Type() {
	case pdata.AttributeValueTypeInt:
		return value.IntVal()
	case pdata.AttributeValueTypeBool:
		return value.BoolVal()
	case pdata.AttributeValueTypeDouble:
		return value.DoubleVal()
	case pdata.AttributeValueTypeString:
		return value.StringVal()
	case pdata.AttributeValueTypeMap:
		values := map[string]interface{}{}
		value.MapVal().Range(func(k string, v pdata.AttributeValue) bool {
			values[k] = attrValue(v)
			return true
		})
		return values
	case pdata.AttributeValueTypeArray:
		arrayVal := value.ArrayVal()
		values := make([]interface{}, arrayVal.Len())
		for i := 0; i < arrayVal.Len(); i++ {
			values[i] = attrValue(arrayVal.At(i))
		}
		return values
	case pdata.AttributeValueTypeNull:
		return nil
	default:
		return nil
	}
}
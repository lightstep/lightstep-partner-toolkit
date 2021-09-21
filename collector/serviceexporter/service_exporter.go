package serviceexporter

import (
	"context"
	"encoding/json"
	"fmt"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/model/pdata"
	"go.uber.org/zap"
	"net/http"
	"strings"
	"sync"
	"time"
)


type ServiceRelationship struct {
	From          string `json:"from"`
	To            string `json:"to"`
	ToOperation string `json:"to_operation"`
	FromOperation string `json:"from_operation"`
	SpanCount     uint64 `json:"span_count"`
	LastSeen      string `json:"last_seen"`
}

//type ServiceResourcesCol struct {
//	Attributes  `json:"attributes"`
//}

type ServiceResourceAttributes map[string]pdata.AttributeMap

func attrsHashString(m pdata.AttributeMap) string {
	return fmt.Sprintf("%v", attrsValue(m))
}

type ServiceResourceCollection struct {
	Attributes map[string]interface{} `json:"attributes"`
	LastSeen string `json:"last_seen"`
}

type ServiceResources struct  {
	Resources     []ServiceResourceCollection `json:"resources"`
	Relationships []ServiceRelationship `json:"relationships"`
}

type ServiceResourceResponse struct {
	Data *ServiceResources `json:"data"`
}

type serviceExporter struct {
	logger              *zap.Logger
	mutex               sync.Mutex
	server              *http.Server
	config              *Config
	serviceResources    *ServiceResources
	spanIdToServiceName map[string]string
	spanIdToOperationName map[string]string
	spanIdToSpanName map[string]string
	relationshipMap map[string]uint64 // key: svc1>svc2, val: ingress operation svc2
	resourceMap map[string]ServiceResourceAttributes
}

func NewServiceExporter(logger *zap.Logger, oCfg *Config) *serviceExporter {
	return &serviceExporter{
		serviceResources: &ServiceResources{
			Resources: make([]ServiceResourceCollection, 0),
		},
		config:              oCfg,
		spanIdToServiceName: make(map[string]string),
		spanIdToOperationName: make(map[string]string),
		relationshipMap: make(map[string]uint64),
		resourceMap: make(map[string]ServiceResourceAttributes),
		logger:              logger,
	}
}

func (e *serviceExporter) Capabilities() consumer.Capabilities {
	return consumer.Capabilities{MutatesData: false}
}

func (e *serviceExporter) addRelationship(parentService string, parentOp string, childService string, operationName string) {
	keyName := fmt.Sprintf("%s%s%s%s%s%s%s", parentService, RelationshipSeparator,
		parentOp, RelationshipSeparator, childService, RelationshipSeparator, operationName)
	_, exists := e.relationshipMap[keyName]; if !exists {
		e.relationshipMap[keyName] = 0
	} else {
		e.relationshipMap[keyName] = e.relationshipMap[keyName] + 1
	}
}

func (e *serviceExporter) ConsumeTraces(_ context.Context, td pdata.Traces) error {
	rss := td.ResourceSpans()
	for i := 0; i < rss.Len(); i++ {
		rs := rss.At(i)
		serviceNameAttr, serviceOk := rs.Resource().Attributes().Get("service.name")
		if !serviceOk {
			serviceNameAttr, serviceOk = rs.Resource().Attributes().Get("lightstep.component_name")
		}
		serviceNameStr := serviceNameAttr.StringVal()

		if serviceOk {
			resourceAttrs := rs.Resource().Attributes()
			hashKey := attrsHashString(resourceAttrs)
			_, ok := e.resourceMap[serviceNameStr]; if !ok {
				e.resourceMap[serviceNameStr] = make(ServiceResourceAttributes)
			}
			_, ok = e.resourceMap[serviceNameStr][hashKey]; if !ok {
				e.resourceMap[serviceNameStr][hashKey] = resourceAttrs
			}
		}
		ils := rs.InstrumentationLibrarySpans()
		for j := 0; j < ils.Len(); j++ {
			is := ils.At(j)
			spans := is.Spans()
			for k := 0; k < spans.Len(); k++ {
				span := spans.At(k)
				if !serviceOk {
					continue
				}

				parentService, parentOk := e.spanIdToServiceName[span.ParentSpanID().HexString()]; if parentOk {
					if parentService != serviceNameStr {
						parentOperation, parentOpOk := e.spanIdToOperationName[span.ParentSpanID().HexString()]; if parentOpOk {
							e.addRelationship(parentService, parentOperation, serviceNameStr, span.Name())
						}
					}
				} else if span.ParentSpanID().Bytes() == [8]byte{0x0} {
					e.addRelationship("__ROOT__", "__ROOT__", serviceNameStr, span.Name())
				}
				e.spanIdToOperationName[span.SpanID().HexString()] = span.Name()
				e.spanIdToServiceName[span.SpanID().HexString()] = serviceNameStr
			}
		}
	}
	return nil
}

const RelationshipSeparator = "\t"

func (e *serviceExporter) Start(_ context.Context, host component.Host) error {
	handler := http.NewServeMux()

	handler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		e.mutex.Lock()
		defer e.mutex.Unlock()

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)

		// rebuild relationship map on each request, probably bad
		e.serviceResources.Relationships = make([]ServiceRelationship, 0)
		for serviceRel, count := range e.relationshipMap {
			services := strings.Split(serviceRel, RelationshipSeparator)
			e.serviceResources.Relationships = append(e.serviceResources.Relationships,
				ServiceRelationship{From: services[0], FromOperation: services[1], To: services[2], SpanCount: count, ToOperation: services[3], LastSeen: time.Now().Format(time.RFC3339)})
		}

		e.serviceResources.Resources = make([]ServiceResourceCollection, 0)
		for _, resourceAttrMap := range e.resourceMap {
			for _, attrs := range resourceAttrMap {
				e.serviceResources.Resources = append(e.serviceResources.Resources,
					ServiceResourceCollection{Attributes: attrsValue(attrs), LastSeen: time.Now().Format(time.RFC3339)})
			}
		}

		data, _ := json.Marshal(ServiceResourceResponse{e.serviceResources})
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

	// backwards compat for pre-opentelemetry service names
	componentName, hasLightstepComponentName := out["lightstep.component_name"]
	_, hasServiceName := out["service.name"]
	if !hasServiceName && hasLightstepComponentName {
		out["service.name"] = componentName
	}
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
	case pdata.AttributeValueTypeEmpty:
		return nil
	default:
		return nil
	}
}
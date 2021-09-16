### OpenTelemetry Collector Experiments

Experimental OpenTelemetry collectors that annotate metrics and traces with external events like PagerDuty incidents, deployments, or chaos experiments. 


### synthetic data instructions

This generates synthetic trace data inside the collector and sends to Lightstep.

```
$ export LS_ACCESS_TOKEN=your token

# to send traces/metrics to public sats (default)
$ docker run -e LS_ACCESS_TOKEN --rm ghcr.io/lightstep/lightstep-partner-toolkit-collector:latest

# to send traces/metrics elsewhere (another collector, Lightstep sats, etc)

# optional: set gRPC transport to insecure (default: `true`, if using dev mode or non-TLS sats)
$ export OTLP_INSECURE=false
# export OTEL_EXPORTER_OTLP_TRACES_ENDPOINT=another-collector:55680

$ docker run -e LS_ACCESS_TOKEN -e OTEL_INSECURE -e OTEL_EXPORTER_OTLP_TRACES_ENDPOINT --rm ghcr.io/lightstep/lightstep-partner-toolkit-collector:latest

```

### webhook demo instructions

This will run the collector and creates webhooks to your local machine via [localtunnel](https://theboroer.github.io/localtunnel-www/).

```
  $ export LS_ACCESS_TOKEN=<your token>
  $ docker-compose up
```

From the output of docker compose, find the metric and trace webhook URLs and add to PagerDuty, Gremlin, or GitHub.

```
  # example output from docker-compose up
  localtunnel-trace-webhook_1   | your url is: https://popular-robin-4.loca.lt
  localtunnel-metric-webhook_1  | your url is: https://short-deer-97.loca.lt
```

You can also manually add and delete attributes that will appear as tags in metrics or traces. Note that the metric pipeline and trace pipeline have separate URLs.

```
    $ curl http://<<hostname>>/upsert?key=foo&value=bar # adds foo=bar to all traces or metrics
    $ curl http://<<hostname>>/delete?key=foo # removes foo to all traces or metrics
```

### running collector (locally)

for log-only output without any processors:
```
  $ docker run --rm -p 55680:55680 -p 55681:55681 \
    -v $(pwd)/config/log-only-config.yml:/etc/otel/config.yaml \
     ghcr.io/lightstep/lightstep-partner-toolkit-collector:latest
```

to send to Lightstep:
```
  $ docker run --rm -p 55680:55680 -p 55681:55681 \
    -v $(pwd)/config/collector-config.yml:/etc/otel/config.yaml \
     ghcr.io/lightstep/lightstep-partner-toolkit-collector:latest
```

send a fake span to the http endpoint:

```
curl -ivX POST -H "Content-Type: application/json" localhost:55681/v1/traces -d '{"resourceSpans":[{"resource":{},"instrumentationLibrarySpans":[{"instrumentationLibrary":{},"spans":[{"traceId":"5b8efff798038103d269b633813fc60c","spanId":"eee19b7ec3c1b173","parentSpanId":"","name":"testSpan","startTimeUnixNano":"1544712660000000000","endTimeUnixNano":"1544712661000000000","attributes":[{"key":"attr1","value":{"intValue":"55"}}],"status":{}}]}]}]}'
```

send a fake metric to the http endpoint:

```
curl -ivX POST "Content-type: application/json" localhost:55681/v1/metrics -d '{"resourceMetrics":[{"resource":{"attributes":[{"key":"service.name","value":{"stringValue":"unknown_service"}}],"droppedAttributesCount":0},"instrumentationLibraryMetrics":[{"metrics":[{"name":"random_count","description":"","unit":"1","doubleSum":{"dataPoints":[{"labels":[{"key":"hostname","value":"test.local"}],"value":36,"startTimeUnixNano":1623690881701000000,"timeUnixNano":1623690893726877700}],"isMonotonic":true,"aggregationTemporality":2}}],"instrumentationLibrary":{"name":"handmade"}}]}]}'
```

### building collector (Docker)

This creates a docker image of the collector with the configuration file `config/collector-config.yml`. By default, it sends OpenTelemetry metrics and traces to Lightstep.

```
  $ docker build -t lightstep/lightstep-partner-toolkit-collector:latest .
```

### building collector (local)

requires go 1.15+

```
  $ go get github.com/open-telemetry/opentelemetry-collector-builder@v0.30.0
  $ opentelemetry-collector-builder --config $(pwd)/builder-config.yml
  $ /tmp/ls-partner-col-distribution/lightstep-partner-collector --config $(pwd)/config/collector-config.yml
```
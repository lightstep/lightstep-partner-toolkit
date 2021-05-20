### OpenTelemetry Collector Experiments

Experimental OpenTelemetry collectors that annotate metrics and traces with external events like PagerDuty incidents, deployments, or chaos experiments. 

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

### building collector (Docker)

This creates a docker image of the collector with the configuration file `config/collector-config.yml`. By default, it sends OpenTelemetry metrics and traces to Lightstep.

```
  $ docker build -t lightstep/lightstep-partner-toolkit-collector:latest .
```

### building collector (local)

requires go 1.15+

```
  $ go get github.com/open-telemetry/opentelemetry-collector-builder
  $ opentelemetry-collector-builder --config $(pwd)/builder-config.yml
  $ /tmp/ls-partner-col-distribution/lightstep-partner-collector --config $(pwd)/config/collector-config.yml
```
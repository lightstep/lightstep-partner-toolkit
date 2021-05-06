### Istio + Lightstep Metrics with the Prometheus Sidecar

Work-in-progress.

#### Adding the Lightstep Prometheus sidecar to Istio

This assume a Istio environment similar to the one that you create in ["Getting Started: Install Istio"](https://istio.io/latest/docs/setup/getting-started/) with istioctl:

##### Create a demo cluster with telemetry + tracing (optional)
```
# You can skip this step if you already have an Istio cluster
$ export LS_ACCESS_TOKEN=<your Lightstep token>
$ istioctl install --set profile=demo \
    --set values.pilot.traceSampling=100 \
    --set values.global.proxy.tracer="lightstep" \
    --set values.global.tracer.lightstep.address="ingest.lightstep.com:443" \
    --set values.global.tracer.lightstep.accessToken=$LS_ACCESS_TOKEN
```

##### Update existing cluster to forward Prometheus data to Lightstep

Next, update Istio to send data to Lightstep with the [OpenTelemetry Prometheus sidecar](https://github.com/lightstep/opentelemetry-prometheus-sidecar):

```
  # 1/ Edit *.yaml file in this directory to set your access token.
  $ vi lightstep-prom-sidecar.yaml

  # 2/ Confirm that prometheus is deployed in the istio-system namespace
  $ kubectl describe deployment prometheus --namespace istio-system
  
  # 3/ Replace Istio Prometheus with Prometheus that runs Lightstep's sidecar
  $ kubectl delete deployment prometheus --namespace istio-system
  $ kubectl apply -f lightstep-prom-sidecar.yaml
```


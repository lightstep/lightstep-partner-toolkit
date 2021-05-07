### Istio + Lightstep Metrics with the Prometheus Sidecar

Work-in-progress.

#### Adding the Lightstep Prometheus sidecar to Istio

This assume a Istio environment similar to the one that you create in ["Getting Started: Install Istio"](https://istio.io/latest/docs/setup/getting-started/) with istioctl:

##### Configure Istio with telemetry + tracing
```
# You can skip this step if you already have an Istio cluster
$ export LS_ACCESS_TOKEN=<your Lightstep token>

$ kubectl create secret generic lightstep.cacert --from-file=cacert.pem

$ kubectl create namespace istio-system

$ istioctl manifest generate -f istio-config.yaml > my-manifest.yaml

$ ... edit manifest to mount volumes

$ kubectl apply -f my-manifest.yaml
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


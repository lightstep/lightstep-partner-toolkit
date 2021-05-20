### Istio + Lightstep Metrics with the Prometheus Sidecar

Work-in-progress.

#### Adding the Lightstep Prometheus sidecar to Istio

This assume a Istio environment similar to the one that you create in ["Getting Started: Install Istio"](https://istio.io/latest/docs/setup/getting-started/) with istioctl:

##### Configure Istio with telemetry + tracing

```
# Edit istio-config-sat.yaml to point to your Microsatellite: https://docs.lightstep.com/docs/install-and-configure-micro-satellites
# Quickstart microsat: run `docker-compose up` in this directory after inserting your token.

# Direct install that points to local microsat
$ istioctl install -f istio-config-sat.yaml

# or: generate and install from yaml
$ istioctl manifest generate -f istio-config-sat.yaml > my-manifest.yaml
$ kubectl create namespace istio-system
$ kubectl apply -f my-manifest.yaml

# Installing the demo app
$ kubectl label namespace default istio-injection=enabled
$ kubectl apply -f https://raw.githubusercontent.com/istio/istio/release-1.9/samples/bookinfo/platform/kube/bookinfo.yaml
$ kubectl apply -f https://raw.githubusercontent.com/istio/istio/release-1.9/samples/bookinfo/networking/bookinfo-gateway.yaml

# Make some traces!
$ curl http://localhost/productpage # on Docker for Mac
```

##### Configure Istio with telemetry + tracing (public sats)

```
$ kubectl create namespace istio-system
$ kubectl create secret generic lightstep.cacert --from-file=cacert.pem --namespace istio-system
$ istioctl manifest generate -f istio-config-sat.yaml > my-manifest.yaml
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



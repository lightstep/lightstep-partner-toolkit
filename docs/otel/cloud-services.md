# Cloud Services and Infrastructure

This guide has three tracks:

| Track | Difficulty |
| ----- | ----- |
| [Use existing telemetry from your service](#) | 🛠 |
| [Add new telemetry to your service](#) | 🛠 🛠 🛠 |
| [Add new telemetry to your SDK or library](#) | 🛠 🛠 🛠 |

## Use Existing Telemetry from your Product or Service

### Design

1. Determine existing metrics, traces, and logs that are available via API.
2. Determine existing metric or tracing formats that are already supported by your product (i.e. OpenTracing, Prometheus exporters).
3. Review the [OpenTelemetry collector](https://github.com/open-telemetry/opentelemetry-collector) and [-contrib](https://github.com/open-telemetry/opentelemetry-collector-contrib) respositories: does an integration already exists, or is it possible to use an existing integration (i.e. prometheus sidecar if your product already outputs prometheus metrics)?
4. If no prior work exists, use the [OpenTelemetry Collector Builder](https://github.com/open-telemetry/opentelemetry-collector-builder) to start a new golang project that pulls data from your service and converts it into an OpenTelemetry-compatible metric, log, or trace.

### Build

* Send logs, metrics or traces telemetry to an OpenTelemetry collector that's configured to ingest data from your service using code you wrote or an existing integration.

### Run

* Run an OpenTelemetry collector that outputs to the console and verify output.

### Example Integrations

| Integration | Description |
| --- | --- |
| [AWS Distro for OpenTelemetry](https://aws.amazon.com/otel) | OpenTelemetry distribution for multiple AWS services. |
| [Lightstep Prometheus Sidecar](https://github.com/lightstep/opentelemetry-prometheus-sidecar) | Convert existing Prometheus metrics into OpenTelemetry metrics. |
| [donut shop demo (EKS)](../../examples/aws) | Example AWS Kubernetes environment that forwards Prometheus metrics to an OpenTelemetry collector. |

## Add New Instrumention to your Product or Service

### Design

* Determine if there are OpenTelemetry SDK(s) available for the language(s) your product or solution is written in.
* Decide what type(s) of telemetry to add to the product and key transactions and metrics that you want to measure first.

### Instrument

* Start coding! Import the language-specific OpenTelemetry API into your product and use it to generate metrics, logs, and traces.

### Run

* Send data from your product to an OpenTelemetry collector that outputs to the console and verify output.

### Example Integrations

| Integration | Description |
| --- | --- |
| [CockroachDB](./examples/cockroachdb) | Instructions for using CockroachDB's native OpenTracing support with Lightstep. |
| [nginx](../../examples/nginx) | Instructions for instrumenting nginx with OpenTelemetry. |
| [Jenkins X](https://github.com/jenkinsci/opentelemetry-plugin) | Publish Jenkins performance metrics and traces to an OpenTelemetry endpoint |
| [Ambassador k8s Initializer](https://lightstep.com/blog/lightstep-and-ambassador/) | Automatically configure a Kubernetes cluster to emit traces using Ambassador's k8s initializer. |

## Add New Instrumention to your SDK or library

### Design

* Determine if there are OpenTelemetry SDK(s) available for the language(s) your SDK or library is written in.

### Instrument

* Start coding! Import the language-specific OpenTelemetry API into your product and use it to generate metrics, logs, and traces.

### Run

* Send data from your product to an OpenTelemetry collector that outputs to the console and verify output.

### Example Integrations

| Instrumentation Package | Description |
| --- | --- |
| [AWS SDK Instrumentation (Node.js)](https://github.com/aspecto-io/opentelemetry-ext-js) | Node.js AWS SDK Instrumentation from [Aspecto](https://github.com/aspecto-io). |
| [MongoDB Python SDK (pymongo)](https://opentelemetry-python-contrib.readthedocs.io/en/latest/instrumentation/pymongo/pymongo.html) | pymongo library instrumentation (open-source). |
| [MySQL Instrumentation](https://github.com/open-telemetry/opentelemetry-js-contrib/tree/main/plugins/node/opentelemetry-instrumentation-mysql) | MySQL client instrumentation for Javascript/Node.js. |
# Cloud Services and Infrastructure

| What are you looking to do? | 
| ----- |
| [Use existing telemetry from your service](#) |
| [Add new telemetry to your service](#add-new-instrumention-to-your-product-or-service) |
| [Add new telemetry to your SDK or library](#add-new-instrumention-to-your-sdk-or-library) |
| [See example integrations](#example-integrations) |

<br/>

## Use Existing Telemetry from your Product or Service

_The goal here is to give your team a unified data pipeline where your cloud services and infrastructure are "speaking the same language" so that your teams and tools ll see the same data._

### Configure Collector

1. Determine existing metrics, traces, and logs that are available via API.
2. Determine existing metric or tracing formats that are already supported by your product (i.e. OpenTracing, Prometheus exporters).
3. Review the [OpenTelemetry collector](https://github.com/open-telemetry/opentelemetry-collector) and [-contrib](https://github.com/open-telemetry/opentelemetry-collector-contrib) respositories: does an integration already exists, or is it possible to use an existing integration (i.e. prometheus sidecar if your product already outputs prometheus metrics)?
4. If no prior work exists, use the [OpenTelemetry Collector Builder](https://github.com/open-telemetry/opentelemetry-collector-builder) to start a new golang project that pulls data from your service and converts it into an OpenTelemetry-compatible metric, log, or trace.

### Pull Data

* Via API, pull in relevant logs, metrics or traces telemetry to an OpenTelemetry collector _receiver_ that's configured to ingest data from your service.
  * Example: [OpenTelemetry Collector Recievers](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/v0.27.0/receiver)

### Run and Verify

1. Send data from your product to an OpenTelemetry collector that outputs to the console and verify output.
2. Optional: verify in an OpenTelemetry production tool of your choice

### Contribute your integration to the OpenTelemetry ecosystem

Make your code usable to as many people as possible! If you're looking for help here, contact us at partnerships@lightstep.com. We'd love to help support you!

<br/>

## Add New Instrumention to your Product or Service

### Instrument

1. **Find the OpenTelemetry SDK for the [language(s)](https://opentelemetry.io/) used by your solution (i.e. nginx is written in C++)**

2. **Import the language-specific OpenTelemetry API and start creating metrics, logs, or traces**
    * Example: [Jenkins X creating OpenTelemetry metrics related to JVM garbage collection](https://github.com/jenkinsci/opentelemetry-plugin/blob/7a6753976df2ca6f5b4b4e4e87772b9e26d6b3db/src/main/java/io/jenkins/plugins/opentelemetry/opentelemetry/instrumentation/runtimemetrics/MemoryPools.java#L67-L77)


### Run and Verify

1. Send data from your product to an OpenTelemetry collector that outputs to the console and verify output.
2. Optional: verify in an OpenTelemetry production tool of your choice

## Add New Instrumention to your SDK or library

### Instrument

1. **Find the OpenTelemetry SDK for the [language(s) or framework(s)](https://opentelemetry.io/) used by your SDK**

2. **Import the language-specific OpenTelemetry API and patch your library method(s)**

    * Example: [Aspecto AWS SDK Instrumentation](https://github.com/aspecto-io/opentelemetry-ext-js/blob/cf5bead74580c520740560c6bd7ca05fc276168c/packages/instrumentation-aws-sdk/src/aws-sdk.ts#L113)

3. **Use the OpenTelemetry docs to add spans, metrics, and logs to annotate requests to your cloud service from customer code with more actionable context**
    * Example: [Aspecto AWS SDK Instrumentation that adds AWS Regions and other metadata](https://github.com/aspecto-io/opentelemetry-ext-js/blob/cf5bead74580c520740560c6bd7ca05fc276168c/packages/instrumentation-aws-sdk/src/aws-sdk.ts#L175-L186)

### Run and Verify

1. Send data from your product to an OpenTelemetry collector that outputs to the console and verify output.
2. Optional: verify in an OpenTelemetry production tool of your choice

<br/>

## Example Integrations

###  Native Instrumentation

| Integration | Description |
| --- | --- |
| [CockroachDB](./examples/cockroachdb) | Instructions for using CockroachDB's native OpenTracing support with Lightstep. |
| [nginx](../../examples/nginx) | Instructions for instrumenting nginx with OpenTelemetry. |
| [Jenkins X](https://github.com/jenkinsci/opentelemetry-plugin) | Publish Jenkins performance metrics and traces to an OpenTelemetry endpoint |
| [Ambassador k8s Initializer](https://lightstep.com/blog/lightstep-and-ambassador/) | Automatically configure a Kubernetes cluster to emit traces using Ambassador's k8s initializer. |

### Collector-based Integration

| Integration | Description |
| --- | --- |
| [AWS Distro for OpenTelemetry](https://aws.amazon.com/otel) | OpenTelemetry distribution for multiple AWS services. |
| [Lightstep Prometheus Sidecar](https://github.com/lightstep/opentelemetry-prometheus-sidecar) | Convert existing Prometheus metrics into OpenTelemetry metrics. |
| [donut shop demo (EKS)](../../examples/aws) | Example AWS Kubernetes environment that forwards Prometheus metrics to an OpenTelemetry collector. |

### SDK or Library Integration

| Instrumentation Package | Description |
| --- | --- |
| [AWS SDK Instrumentation (Node.js)](https://github.com/aspecto-io/opentelemetry-ext-js) | Node.js AWS SDK Instrumentation from [Aspecto](https://github.com/aspecto-io). |
| [MongoDB Python SDK (pymongo)](https://opentelemetry-python-contrib.readthedocs.io/en/latest/instrumentation/pymongo/pymongo.html) | pymongo library instrumentation (open-source). |
| [MySQL Instrumentation](https://github.com/open-telemetry/opentelemetry-js-contrib/tree/main/plugins/node/opentelemetry-instrumentation-mysql) | MySQL client instrumentation for Javascript/Node.js. |
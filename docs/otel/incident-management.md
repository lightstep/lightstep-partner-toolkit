# Incident Management

This guide has three tracks:

| Track | Difficulty |
| ----- | ----- |
| [Use existing telemetry from your service](#) | 🛠🛠 |

## Use Existing Telemetry from your Product or Service

### Design

1. Determine existing event(s) available from your product that can fire a webhook.
2. Determine existing metric or tracing formats that are already supported by your product (i.e. OpenTracing, Prometheus exporters).
3. Review the [OpenTelemetry collector](https://github.com/open-telemetry/opentelemetry-collector) and [-contrib](https://github.com/open-telemetry/opentelemetry-collector-contrib) respositories: does an integration already exists, or is it possible to use an existing integration?
4. If no prior work exists, use the [OpenTelemetry Collector Builder](https://github.com/open-telemetry/opentelemetry-collector-builder) to start a new golang project that receives events from your service and annotates relevant OpenTelemetry metrics, logs, and traces.

### Build

* Send logs, metrics or traces telemetry to an OpenTelemetry collector that's configured to recieve webhooks from your service using code you wrote or an existing integration.

### Run

* Run an OpenTelemetry collector that outputs to the console and verify output.

### Example Integrations

| Integration | Description |
| --- | --- |
| [webhookprocessor](../../collector) | OpenTelemetry collector that can receive PagerDuty incident webhooks. |

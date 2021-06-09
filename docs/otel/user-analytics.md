# User Analytics

This guide has three tracks:

| Track | Difficulty |
| ----- | ----- |
| [Add new telemetry to your analytics SDK or library](#) | 🛠 🛠 🛠 |

## Add New Instrumention to your SDK or library

### Design

* Determine if there are OpenTelemetry SDK(s) available for the language(s) your SDK or library is written in.

### Instrument

* Start coding! Import the language-specific OpenTelemetry API into your product and use it to generate metrics, logs, and traces.

### Run

* Send data from your product to an OpenTelemetry collector that outputs to the console and verify output.

### Example Integrations

> These example solutions use plugins to generate metrics, logs or traces by automatically-instrumenting Javascript libaries or frameworks. No change to the underlying library or framework is needed.

| Instrumentation Package | Instrumented Package |
| --- | --- |
| [opentelemetry-plugin-segment-node](../../js/packages/opentelemetry-plugin-segment-node) | [`analytics-node`](https://github.com/segmentio/analytics-node) |

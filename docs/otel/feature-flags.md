# Feature Flags

This guide has three tracks:

| Track | Difficulty |
| ----- | ----- |
| [Add new telemetry to your feature flag SDK or library](#) | 🛠 🛠 🛠 |

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
| [opentelemetry-plugin-splitio](../../js/packages/opentelemetry-plugin-splitio) | [`@splitsoftware/splitio`](https://github.com/splitio/javascript-client) |
| [opentelemetry-plugin-launchdarkly-node-server](../../js/packages/opentelemetry-plugin-launchdarkly-node-server) | [`launchdarkly-node-server-sdk`](https://github.com/launchdarkly/node-server-sdk) |
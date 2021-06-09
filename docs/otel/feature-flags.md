# Feature Flags

Adding OpenTelemetry spans to your feature flags implementation provides a **clear connection between the code, your product and engineering teams, and what your customers are seeing in production**. The full context from OpenTelemetry removes the mystery when managing your deployments and customer experience.

| What are you looking to do? | 
| ----- | 
| [Add new telemetry to your feature flag SDK or library](#add-new-instrumention-to-your-sdk-or-library) |
| [See example integrations](#example-integrations) |

## Add New Instrumention to your SDK or library

### Instrument

1. Determine if there are OpenTelemetry SDK(s) available for the language(s) your SDK or library is written in.
2. Start coding! Import the language-specific OpenTelemetry API into your product and use it to generate metrics, logs, and traces.

### Run

* Send data from your product to an OpenTelemetry collector that outputs to the console and verify output.

<br/>

## Example Integrations

> These example solutions use plugins to generate metrics, logs or traces by automatically-instrumenting Javascript libaries or frameworks. No change to the underlying library or framework is needed.

| Instrumentation Package | Instrumented Package |
| --- | --- |
| [opentelemetry-plugin-splitio](../../js/packages/opentelemetry-plugin-splitio) | [`@splitsoftware/splitio`](https://github.com/splitio/javascript-client) |
| [opentelemetry-plugin-launchdarkly-node-server](../../js/packages/opentelemetry-plugin-launchdarkly-node-server) | [`launchdarkly-node-server-sdk`](https://github.com/launchdarkly/node-server-sdk) |

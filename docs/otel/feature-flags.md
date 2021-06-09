# Feature Flags

Connecting OpenTelemetry spans to your feature flags implementation provides a **clear connection between the code, your product and engineering teams, and what your customers are seeing in production**. The full context from OpenTelemetry removes the mystery when managing your deployments and customer experience.

| What are you looking to do? | 
| ----- | 
| [Add new telemetry to your feature flag SDK or library](#add-new-instrumention-to-your-sdk-or-library) |
| [See example integrations](#example-integrations) |

## Add New Instrumention to your SDK or library

_The goal here is to connect feature flag state and customer experience to the full picture of both production system state so that engineers can move faster in building the product._

### Instrument

1. **Find the OpenTelemetry SDK for the [language(s) or framework(s)](https://opentelemetry.io/) used by your SDK**

2. **Import the language-specific OpenTelemetry API and patch your library method(s)**
   * [Example from JavaScript](https://github.com/lightstep/lightstep-partner-toolkit/blob/main/js/packages/opentelemetry-plugin-rollbar/src/rollbar.ts#L1)

3. **Use the OpenTelemetry docs to add spans, metrics, and logs to annotate calls to retrieve feature flags produced by your SDK with more actionable context**
   * As an example, here the Rollbar SDK creates [an attribute to linking back to the specific error](https://github.com/lightstep/lightstep-partner-toolkit/blob/d42c616a227dedbc013e698bdee454f4844d571c/js/packages/opentelemetry-plugin-rollbar/src/rollbar.ts#L48) so the developer can link to the full context of the error

### Run and Verify

1. Send data from your product to an OpenTelemetry collector that outputs to the console and verify output.
2. Optional: verify in an OpenTelemetry production tool of your choice

### Contribute your integration to the OpenTelemetry ecosystem

Make your code usable to as many people as possible! If you're looking for help here, contact us at partnerships@lightstep.com. We'd love to help support you!

<br/>

## Example Integrations

> These example solutions use plugins to generate metrics, logs or traces by automatically-instrumenting Javascript libaries or frameworks. No change to the underlying library or framework is needed.

| Instrumentation Package | Instrumented Package |
| --- | --- |
| [opentelemetry-plugin-splitio](../../js/packages/opentelemetry-plugin-splitio) | [`@splitsoftware/splitio`](https://github.com/splitio/javascript-client) |
| [opentelemetry-plugin-launchdarkly-node-server](../../js/packages/opentelemetry-plugin-launchdarkly-node-server) | [`launchdarkly-node-server-sdk`](https://github.com/launchdarkly/node-server-sdk) |

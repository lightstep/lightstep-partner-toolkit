# Errors

Connecting OpenTelemetry spans to your SDK or library's errors give the context to connect an error across teams, tools, and products to allow developers to resolve issues faster.  Example such as adding spans with **project IDs**, **URLs**, or **segmentation information** can give developers directly actionable information or direct them to who to pull into the conversation.

| What are you looking to do? | 
| ----- | 
| [Add OpenTelemetry support for errors emitted by your SDK or library](#add-opentelemetry-support-for-errors-emitted-by-your-sdk-or-library) |
| [See example integrations](#example-integrations) |

<br/>

## Add OpenTelemetry support for errors emitted by your SDK or library

_The goal here is to connect error analysis to the full view of the distributed system by linking error events to the current tracing context._


### Instrument

1. **Find the OpenTelemetry SDK for the [language(s) or framework(s)](https://opentelemetry.io/) used by your SDK**

2. **Import the language-specific OpenTelemetry API and patch your library method(s)**
   * [Example from JavaScript](https://github.com/lightstep/lightstep-partner-toolkit/blob/main/js/packages/opentelemetry-plugin-rollbar/src/rollbar.ts#L1)

3. **Use the OpenTelemetry docs to add spans, metrics, and logs to annotate errors produced by your SDK with more actionable context**
   * As an example, here the Rollbar SDK creates [an attribute to linking back to the specific error](https://github.com/lightstep/lightstep-partner-toolkit/blob/d42c616a227dedbc013e698bdee454f4844d571c/js/packages/opentelemetry-plugin-rollbar/src/rollbar.ts#L48) so the developer can link to the full context of the error

### Run and Verify

1. Send data from your product to an OpenTelemetry collector that outputs to the console and verify output.
2. Optional: verify in an OpenTelemetry production tool of your choice

### Contribute your integration to the OpenTelemetry ecosystem

Make your code usable to as many people as possible! If you're looking for help here, contact us at partnerships@lightstep.com. We'd love to help support you!

<br/>

## Example Integrations

> These example solutions use plugins to generate metrics, logs or traces by automatically-instrumenting Javascript libaries or frameworks. No change to the underlying library or framework is needed.

| Instrumentation Package | Instrumented Package | What it does |
| --- | --- | --- |
| [opentelemetry-plugin-rollbar](./js/packages/opentelemetry-plugin-rollbar) | [`rollbar`](https://github.com/rollbar/rollbar.js/) | emits OpenTelemetry data from the [Rollbar](https://rollbar.com/) SDK to connect Rollbar error analysis to OpenTelemetry trace data|

# User Analytics

Connecting OpenTelemetry spans to your analytics SDK or library connects end user and product data back to the  underlying production system data. This creates a unified data pipeline allowing your product and engineering teams to work more readily in sync. Example such as adding spans with **customer IDs**, **account type**, or **segmentation information** can give developers directly actionable information that bridges technical performance with customer experience.

| What are you looking to do? | 
| ----- |
| [Add new telemetry to your analytics SDK or library](#add-new-instrumention-to-your-analytics-sdk-or-library) |
| [See example integrations](#example-integrations) |

<br/>

## Add New Instrumention to your analytics SDK or library

_The goal here is to bridge the gap between a product-centric view of user analytics data and the engineering-centric view of the production system data, leading to a more coordinated team and better product._

### Instrument

1. **Find the OpenTelemetry SDK for the [language(s) or framework(s)](https://opentelemetry.io/) used by your SDK**

2. **Import the language-specific OpenTelemetry API and patch your library method(s)**
   * [Example from JavaScript](https://github.com/lightstep/lightstep-partner-toolkit/blob/d42c616a227dedbc013e698bdee454f4844d571c/js/packages/opentelemetry-plugin-segment-node/src/segment.ts#L8)

3. **Use the OpenTelemetry docs to add spans, metrics, and logs to annotate analytics calls produced by your SDK with more actionable context**
   * As an example, here the Segement SDK creates [an event added to the span for every `track` event](https://github.com/lightstep/lightstep-partner-toolkit/blob/d42c616a227dedbc013e698bdee454f4844d571c/js/packages/opentelemetry-plugin-segment-node/src/segment.ts#L37) so the developer can see relevant analytics metadata.

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
| [opentelemetry-plugin-segment-node](../../js/packages/opentelemetry-plugin-segment-node) | [`analytics-node`](https://github.com/segmentio/analytics-node) |

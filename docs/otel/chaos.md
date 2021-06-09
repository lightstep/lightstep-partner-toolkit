# Chaos Engineering

Connecting OpenTelemetry spans to ongoing chaos experiments creates **clear links between the code, your product and engineering teams, and an active chaos experiment**. The full context from OpenTelemetry removes the mystery when figuring out "what happened?" during chaos engineering experiments.

| What are you looking to do? | 
| ----- |
| [Use existing telemetry from your service](#use-existing-events-from-your-product-or-service) |
| [See example integrations](#example-integrations) |

<br/>

## Use Existing Events from your Product or Service

_The goal here is to annotate testing data with rich telemetry data so that the developer can quickly isolate issues and improve the system._

### Configure Collector

1. Determine existing event(s) available from your product that can fire a webhook.
2. Determine existing metric or tracing formats that are already supported by your product (i.e. OpenTracing, Prometheus exporters).
3. Review the [OpenTelemetry collector](https://github.com/open-telemetry/opentelemetry-collector) and [-contrib](https://github.com/open-telemetry/opentelemetry-collector-contrib) respositories: does an integration already exists, or is it possible to use an existing integration?
4. If no prior work exists, use the [OpenTelemetry Collector Builder](https://github.com/open-telemetry/opentelemetry-collector-builder) to start a new golang project that receives events from your service and annotates relevant OpenTelemetry metrics, logs, and traces.
    * [Example Collector](https://github.com/lightstep/lightstep-partner-toolkit/tree/main/collector)

### Receive Event(s)

1. Send logs, metrics or traces telemetry to an OpenTelemetry collector that's configured to recieve webhooks from your service using code you wrote or an existing integration. 
2. Using your new custom collector processor, attach the chaos experiment id to the span when the webhook event trigger is fired, and remove it when the webhook resolution trigger is fired.

### Run and Verify

1. Send webhooks from your product to an OpenTelemetry collector configured for console output. Verify outpiut in the console.
2. Optional: verify in an OpenTelemetry production tool of your choice

### Contribute your integration to the OpenTelemetry ecosystem

Make your code usable to as many people as possible! If you're looking for help here, contact us at partnerships@lightstep.com. We'd love to help support you!

<br/>

## Example Integrations

| Integration | Description |
| --- | --- |
| [webhookprocessor](../../collector) | OpenTelemetry collector that can receive Gremlin webhooks. |

# Incident Management

Connecting OpenTelemetry spans to incidents creates **clear links between the code, your product and engineering teams, and active incidents**. The full context from OpenTelemetry removes the mystery when figuring out what technical data is related to ongoing problems or issues tracked in incident management tools, and an invaluable resource for post-mortem analysis.

| What are you looking to do? | 
| ----- |
| [Use existing events from your service](#use-existing-events-from-your-product-or-service) 
| [See example integrations](#example-integrations) |

<br/>

## Use Existing Events from your Product or Service

_The goal here is to produce richer telemetry data for the developer in response to an event or incident occurring._

### Configure Collector

1. Determine existing event(s) available from your product that can fire a webhook, like an incident.  This webhook can trigger the collector to generate additional spans, attributes, metrics, or other data which will help diagnose the incident.

2. Review the [OpenTelemetry collector](https://github.com/open-telemetry/opentelemetry-collector) and [-contrib](https://github.com/open-telemetry/opentelemetry-collector-contrib) respositories: does an integration already exists, or is it possible to use an existing integration?

3. If no prior work exists, use the [OpenTelemetry Collector Builder](https://github.com/open-telemetry/opentelemetry-collector-builder) to start a new golang project that receives events from your service and annotates relevant OpenTelemetry metrics, logs, and traces.

    * [Example Collector](https://github.com/lightstep/lightstep-partner-toolkit/tree/main/collector)

### Receive Event(s)

1. Send logs, metrics or traces telemetry to an OpenTelemetry collector that's configured to receive webhooks from your service using code you wrote or an existing integration. 
2. Using your new custom collector processor, attach the incident id to the span when the webhook trigger is fired, and remove it when the webhook resolution trigger is fired.

### Run and Verify

1. Send webhooks from your product to an OpenTelemetry collector configured for console output. Verify outpiut in the console.
2. Optional: verify in an OpenTelemetry production tool of your choice

### Contribute your integration to the OpenTelemetry ecosystem

Make your code usable to as many people as possible! If you're looking for help here, contact us at partnerships@lightstep.com. We'd love to help support you!

<br/>

## Example Integrations

| Integration | Description |
| --- | --- |
| [webhookprocessor](../../collector) | An OpenTelemetry collector that receives PagerDuty incident webhooks so that telemetry data can be quickly associated to a particular incident |

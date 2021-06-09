### 🪕 Instrument with OpenTelemetry

To get started, you'll need to determine how you want to use the OpenTelemetry APIs. The pattern you choose depends on what kind of solution or product you have: the integration will look different for a database versus a error tracking library, for example.

> 💡 Not sure if OpenTelemetry is relevant for your product, tool, or solution? Here's a quick test: does it contain any data or context that can help people understand their apps or services? This includes anything from technical or performance data to metadata about teams, deploys or incidents. 

There are three main categories of instrumentation for tools or services:

1. If your product, tool or solution is a __service__ like a cloud service, database, firewall or application gateway that customers run or you run for customers (i.e. SaaS): see examples under *Code-based Instrumentation*. You'll use the OpenTelemetry to generate metrics, traces or logs from your product by changing its source code.

2. If your product, tool or solution is a __service__ that _already_ has metrics, logs, traces or events that you want to convert to OpenTelemetry, see *Collector-based Integration*. You'll write an external adapter that converts existing telemetry into OpenTelemetry.

3. If your product, tool, or solution includes a __library or SDK client__ written in Node.js, TypeScript, or Python that customers run alongside their service or application code like a feature flag library or cloud SDK see *Code-based Instrumentation with Plugins*. For other languages, see #1.

> 💡 Before you build, double-check the [OpenTelemetry registry](https://opentelemetry.io/registry/) to see if someone already has contributed code related to your project.

### OpenTelemetry for Browsers

#### Quick start for browser-based integrations

A detailed tutorial is available in Lightstep's [OpenTelemetry examples](https://github.com/lightstep/opentelemetry-examples/tree/main/browser).

To perform basic instrumentation and start seeing data, add this to the top of your HTML:

```
  <!--
    https://www.w3.org/TR/trace-context/
    Set the `traceparent` in the server's HTML template code. It should be
    dynamically generated server side to have the server's request trace Id,
    a parent span Id that was set on the server's request span, and the trace
    flags to indicate the server's sampling decision
    (01 = sampled, 00 = notsampled).
    '{version}-{traceId}-{spanId}-{sampleDecision}'
  -->
  <meta name="traceparent" content="00-ab42124a3c573678d4d8b21ba52df3bf-d21f7bc17caa5aba-01">
  
  <script>window.LS_ACCESS_TOKEN = 'your-lightstep-token'</script>
  <script>window.LS_SERVICE_NAME = 'browser';</script>
  <script src="https://cdn.jsdelivr.net/gh/lightstep/lightstep-partner-sdk/examples/browser/dist/tracer.js"/>
```

... then make some HTTP requests. Data will automatically be send to Lighstep's OpenTelemetry endpoint for collection.
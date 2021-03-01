const { propagation, trace } = require('@opentelemetry/api');
const { NodeTracerProvider } = require('@opentelemetry/node');
const { SimpleSpanProcessor, ConsoleSpanExporter } = require('@opentelemetry/tracing');
const { CollectorTraceExporter } = require('@opentelemetry/exporter-collector');
const { B3Propagator } = require('@opentelemetry/propagator-b3');
const { registerInstrumentations } = require('@opentelemetry/instrumentation');
const { AwsInstrumentation } = require('opentelemetry-instrumentation-aws-sdk');

const path = require('path');

propagation.setGlobalPropagator(new B3Propagator());

module.exports = (serviceName) => {
  const provider = new NodeTracerProvider();

  // This sends data to Lightstep by default
  // Lots of other exporters are supported, see
  // https://opentelemetry.io/registry/
  const exporter = new CollectorTraceExporter({
    serviceName,
    url: 'https://ingest.lightstep.com/traces/otlp/v0.6',
    headers: {
      'Lightstep-Access-Token':
        process.env.LS_ACCESS_TOKEN,
    },
  });

  provider.addSpanProcessor(new SimpleSpanProcessor(new ConsoleSpanExporter()));
  provider.addSpanProcessor(new SimpleSpanProcessor(exporter));

  registerInstrumentations({
    tracerProvider: provider,
    instrumentations: [
      new AwsInstrumentation({
        suppressInternalInstrumentation: true,
        preRequestHook: (span, request) => {
          if (span.attributes['aws.service.api'] === 's3') {
            span.setAttribute('s3.bucket.name', request.params.Bucket);
          }
        },
      }),
      {
        plugins: {
          express: {
            enabled: true,
          },
          rollbar: {
            path: path.join(__dirname, '../../js/packages/opentelemetry-plugin-rollbar/build/src'),
            enabled: true,
          },
          '@splitsoftware/splitio': {
            path: path.join(__dirname, '../../js/packages/opentelemetry-plugin-splitio/build/src'),
            enabled: true,
          },
          'launchdarkly-node-server-sdk': {
            path: path.join(__dirname, '../../js/packages/opentelemetry-plugin-launchdarkly-node-server/build/src'),
            enabled: true,
          },
          'analytics-node': {
            path: path.join(__dirname, '../../js/packages/opentelemetry-plugin-segment-node/build/src'),
            enabled: true,
          },
        },
      },
    ],
  });

  // Initialize the OpenTelemetry APIs
  provider.register();

  return trace.getTracer(serviceName);
};

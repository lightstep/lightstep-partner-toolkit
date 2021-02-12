const opentelemetry = require('@opentelemetry/api');
const { ConsoleLogger, LogLevel } = require('@opentelemetry/core');
const { NodeTracerProvider } = require('@opentelemetry/node');
const { SimpleSpanProcessor } = require('@opentelemetry/tracing');
const { CollectorTraceExporter } = require('@opentelemetry/exporter-collector');
const { B3Propagator } = require('@opentelemetry/propagator-b3');
const path = require('path');

opentelemetry.propagation.setGlobalPropagator(new B3Propagator());

module.exports = (serviceName) => {
  const provider = new NodeTracerProvider({
    plugins: {
      '@splitsoftware/splitio': {
        path: path.join(__dirname, '../node_modules/@lightstep/opentelemetry-plugin-splitio'),
        enabled: true,
      },
      'launchdarkly-node-server-sdk': {
        path: path.join(__dirname, '../node_modules/@lightstep/opentelemetry-plugin-launchdarkly-node-server'),
        enabled: true,
      },
    },
  });

  // This sends data to Lightstep by default
  // Lots of other exporters are supported, see
  // https://opentelemetry.io/registry/
  const exporter = new CollectorTraceExporter({
    serviceName,
    logger: new ConsoleLogger(LogLevel.DEBUG),
    url: 'https://ingest.lightstep.com/traces/otlp/v0.6',
    headers: {
      'Lightstep-Access-Token':
        process.env.LS_ACCESS_TOKEN,
    },
  });

  provider.addSpanProcessor(new SimpleSpanProcessor(exporter));

  // Initialize the OpenTelemetry APIs
  provider.register();

  return opentelemetry.trace.getTracer(serviceName);
};

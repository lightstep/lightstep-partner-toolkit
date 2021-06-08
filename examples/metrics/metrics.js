const { MeterProvider, ConsoleMetricExporter } = require('@opentelemetry/metrics');
const { CollectorMetricExporter } =  require('@opentelemetry/exporter-collector');
const { HostMetrics } = require('@opentelemetry/host-metrics');
const opentelemetry = require('@opentelemetry/api');

opentelemetry.diag.setLogger(
  new opentelemetry.DiagConsoleLogger(), opentelemetry.DiagLogLevel.DEBUG,
);

const collectorOptions = {
  // url is optional and can be omitted - default is localhost:4317
  headers: {
    'lightstep-access-token': process.env.LS_ACCESS_TOKEN
  },
  url: 'https://ingest.lightstep.com:443/metrics/otlp/v0.6',
};
const exporter = new CollectorMetricExporter(collectorOptions); // new ConsoleMetricExporter()

// Register the exporter
const provider = new MeterProvider({
  exporter,
  interval: 1000,
})

const meter = provider.getMeter('example-meter');

// collect host metrics
const hostMetrics = new HostMetrics({ meter, name: 'example-host-metrics' });
hostMetrics.start();

// Now, start recording data
const counter = meter.createCounter('random_count');

setInterval(function() {
  console.log('generating metrics...');
  counter.add(Math.floor(Math.random() * 10), { 'hostname': require("os").hostname() });
}, 1500);

const path = require('path');

const {
  lightstep, opentelemetry,
} = require('lightstep-opentelemetry-launcher-node');

const { AwsInstrumentation } = require('opentelemetry-instrumentation-aws-sdk');

module.exports = async (serviceName) => {
  const sdk = lightstep.configureOpenTelemetry({
    accessToken: process.env.LS_ACCESS_TOKEN,
    serviceName,
    instrumentations: [
      new AwsInstrumentation({
        suppressInternalInstrumentation: true,
        preRequestHook: (span, request) => {
          if (span.attributes['aws.service.api'] === 'S3') {
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
            path: '@lightstep/opentelemetry-plugin-rollbar',
            enabled: true,
          },
          '@splitsoftware/splitio': {
            path: '@lightstep/opentelemetry-plugin-splitio',
            enabled: true,
          },
          'launchdarkly-node-server-sdk': {
            path: '@lightstep/opentelemetry-plugin-launchdarkly-node-server',
            enabled: true,
          },
          'analytics-node': {
            path: '@lightstep/opentelemetry-plugin-segment-node',
            enabled: true,
          },
          /*rollbar: {
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
          },*/
        },
      },
    ],
  });

  // Setting the default Global logger to use the Console
  // And optionally change the logging level (Defaults to INFO)
  opentelemetry.diag.setLogger(
    new opentelemetry.DiagConsoleLogger(), opentelemetry.DiagLogLevel.DEBUG,
  );

  return sdk.start();
};

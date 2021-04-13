# Lightstep Partner Toolkit

Technical resources for Lightstep Partners.

### Getting Started

New to OpenTelemetry, Lightstep, or observability? Check out [technical resources](./resources.md) to learn more.

### OpenTelemetry Instrumentation

[Demo app](./demo/readme.md) and working plugin examples of how to integrate external products with the OpenTelemetry standard.

#### Instrumentation Plugins (Node.js)

Instrumentation has been tested with OpenTelemetry `0.18.0`. Packages are available via GitHub's package registry. 

To install via npm, run:

```
$ npm_config_registry=https://npm.pkg.github.com/lightstep npm install --save <package-name>
```

| Instrumentation Package | Instrumented Package |
| --- | --- |
| [opentelemetry-plugin-splitio](./js/packages/opentelemetry-plugin-splitio) | [`@splitsoftware/splitio`](https://github.com/splitio/javascript-client) |
| [opentelemetry-plugin-launchdarkly-node-server](./js/packages/opentelemetry-plugin-launchdarkly-node-server) | [`launchdarkly-node-server-sdk`](https://github.com/launchdarkly/node-server-sdk) |
| [opentelemetry-plugin-rollbar](./js/packages/opentelemetry-plugin-rollbar) | [`rollbar`](https://github.com/rollbar/rollbar.js/) |
| [opentelemetry-plugin-segment-node](./js/packages/opentelemetry-plugin-segment-node) | [`analytics-node`](https://github.com/segmentio/analytics-node) |

#### OpenTelemetry Collector Processors

Experimental processors for the [OpenTelemetry Collector](https://github.com/open-telemetry/opentelemetry-collector).

| Processor | Description | Partner Integrations |
| --- | --- | --- |
| [webhookprocessor](./collector/webhookprocessor) | Annotates spans with metadata provided by webhooks. | PagerDuty, Gremlin, GitHub Deployments |
| [backstageprocessor](./collector/webhookprocessor) | Annotates spans with service catalog metadata. | [Backstage](https://backstage.io/) |

### Other OpenTelemetry Integrations

| Integration | Description |
| --- | --- |
| [AWS SDK Instrumentation (Node.js)](https://github.com/aspecto-io/opentelemetry-ext-js) | Node.js AWS SDK Instrumentation from [Aspecto](https://github.com/aspecto-io). |
| [CockroachDB](./examples/cockroachdb) | Instructions for using CockroachDB's native OpenTracing support with Lightstep. |
| [nginx](./examples/nginx) | Instructions for instrumenting nginx with OpenTelemetry. |
| [Ambassador k8s Initializer](https://lightstep.com/blog/lightstep-and-ambassador/) | Automatically configure a Kubernetes cluster to emit traces using Ambassador's k8s initializer. |
| [Jenkins X](https://github.com/jenkinsci/opentelemetry-plugin) | Publish Jenkins performance metrics and traces to an OpenTelemetry endpoint |

### Demo

Node.js web application that uses Lightstep partner plugins. See instructions in [`./demo`](./demo) to run the app.




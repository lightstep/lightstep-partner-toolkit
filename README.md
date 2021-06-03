# Lightstep Partner Toolkit

Technical toolkit for Lightstep partners that want to adopt OpenTelemetry. 

Looking for Lightstep's REST API? See the documentation [here](https://api-docs.lightstep.com/reference).

## 💻 Demo

Run [Donut Shop](./demo/readme.md) locally to understand how OpenTelemetry can connect different tools and solutions together. Our example app creates distributed traces that contain data related to relevant feature flags, errors, user analytics, and supporting cloud services.

> 💡 It's also possible to run Donut Shop on a fully-featured Kuberentes cluster using the AWS CDK. See example [here](./examples/aws).

## 📓 How to integrate OpenTelemetry

This toolkit has technical resources and examples you can follow to add OpenTelemetry metrics, logs, or traces to your product or tool. 

> 💡 If you're completely new to OpenTelemetry, check out [technical resources](./resources.md) to learn more.

### 🪕 Instrument

To get started, you'll need to determine how you want to use the OpenTelemetry APIs. The pattern you choose depends on what kind of solution or product you have: the integration will look different for a database versus a error tracking library, for example.

> 💡 Not sure if OpenTelemetry is relevant for your product, tool, or solution? Here's a quick test: does it contain any data or context that can help people understand their apps or services? This includes anything from technical or performance data to metadata about teams, deploys or incidents. 

There are three main categories of instrumentation for tools or services:

1. If your product, tool or solution is a __service__ like a cloud service, database, firewall or application gateway that customers run or you run for customers (i.e. SaaS): see examples under *Code-based Instrumentation*. You'll use the OpenTelemetry to generate metrics, traces or logs from your product by changing its source code.

2. If your product, tool or solution is a __service__ that _already_ has metrics, logs, traces or events that you want to convert to OpenTelemetry, see *Collector-based Integration*. You'll write an external adapter that converts existing telemetry into OpenTelemetry.

3. If your product, tool, or solution includes a __library or SDK client__ written in Node.js, TypeScript, or Python that customers run alongside their service or application code like a feature flag library or cloud SDK see *Code-based Instrumentation with Plugins*. For other languages, see #1.

> 💡 Before you build, double-check the [OpenTelemetry registry](https://opentelemetry.io/registry/) to see if someone already has contributed code related to your project.

#### Code-based Instrumentation

These solutions have directly implemented OpenTelemetry to generate new metrics, logs or traces.

##### Examples

| Integration | Description |
| --- | --- |
| [CockroachDB](./examples/cockroachdb) | Instructions for using CockroachDB's native OpenTracing support with Lightstep. |
| [nginx](./examples/nginx) | Instructions for instrumenting nginx with OpenTelemetry. |
| [Ambassador k8s Initializer](https://lightstep.com/blog/lightstep-and-ambassador/) | Automatically configure a Kubernetes cluster to emit traces using Ambassador's k8s initializer. |
| [Jenkins X](https://github.com/jenkinsci/opentelemetry-plugin) | Publish Jenkins performance metrics and traces to an OpenTelemetry endpoint |

#### Code-based Instrumentation with Plugins (TypeScript/Node.js)

These solutions use plugins to generate metrics, logs or traces by automatically-instrumenting libaries or frameworks. No change to the underlying library or framework is needed.

##### Examples

Below are example OpenTelemetry plugins for Node.js. They can be installed in a node.js project with Github's package registry:
```
$ npm_config_registry=https://npm.pkg.github.com/lightstep npm install --save <package-name>
```

| Instrumentation Package | Instrumented Package |
| --- | --- |
| [opentelemetry-plugin-splitio](./js/packages/opentelemetry-plugin-splitio) | [`@splitsoftware/splitio`](https://github.com/splitio/javascript-client) |
| [opentelemetry-plugin-launchdarkly-node-server](./js/packages/opentelemetry-plugin-launchdarkly-node-server) | [`launchdarkly-node-server-sdk`](https://github.com/launchdarkly/node-server-sdk) |
| [opentelemetry-plugin-rollbar](./js/packages/opentelemetry-plugin-rollbar) | [`rollbar`](https://github.com/rollbar/rollbar.js/) |
| [opentelemetry-plugin-segment-node](./js/packages/opentelemetry-plugin-segment-node) | [`analytics-node`](https://github.com/segmentio/analytics-node) |

##### Other Plugin Examples

| Instrumentation Package | Description |
| --- | --- |
| [AWS SDK Instrumentation (Node.js)](https://github.com/aspecto-io/opentelemetry-ext-js) | Node.js AWS SDK Instrumentation from [Aspecto](https://github.com/aspecto-io). |

#### [Collector](https://github.com/open-telemetry/opentelemetry-collector)-based Integration

##### Processor Examples

| Processor | Description | Partner Integrations |
| --- | --- | --- |
| [webhookprocessor](./collector/webhookprocessor) | Annotates spans with metadata provided by webhooks. | PagerDuty, Gremlin, GitHub Deployments |
| [backstageprocessor](./collector/webhookprocessor) | Annotates spans with service catalog metadata. | [Backstage](https://backstage.io/) |

### 📈 Verify

TBD.

### ↪ Contribute

* Join [CNCF Slack](http://slack.cncf.io/) and say hello on the #opentelemetry channel.
* Find relevant OpenTelemetry [Special Interest Groups](https://github.com/open-telemetry/community#special-interest-groups) and join the discussion or scheduled meeting.
* Open a pull request on projects and contribute your code. 



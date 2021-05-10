# AWS Example Observability Stack

This is a Kubernetes cluster running on AWS instrumented for OpenTelemetry.

Work in progress. This stack is managed using [AWS CDK](https://aws.amazon.com/cdk/).

## Requirements
* `AWS CDK`
* Lightstep
* nginx plus image in private repo (if using nginx plus)

## Deploying

```sh
  # Set LS_ACCESS_TOKEN to send data to Lightstep
  $ export LS_ACCESS_TOKEN=...
  # Set secret for access to private Docker registry with your nginx plus image (optional)
  $ export DOCKER_CONFIG_BASE64=...
  # Install cloud dependencies for CDK (only needed first time)
  $ cdk bookstrap
  # Deploy this stack
  $ cdk deploy
```

## Useful commands

 * `npm run build`   compile typescript to js
 * `npm run watch`   watch for changes and compile
 * `npm run test`    perform the jest unit tests
 * `cdk deploy`      deploy this stack to your default AWS account/region
 * `cdk diff`        compare deployed stack with current state
 * `cdk synth`       emits the synthesized CloudFormation template

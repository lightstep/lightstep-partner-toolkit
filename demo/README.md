###  OpenTelemetry + Lightstep Partners Node.js Demo App

Simple Node.js web app that emits OpenTelemetry with Lightstep partner plugins enabled. To preview this in the cloud, **[Run in Codesandbox](https://codesandbox.io/s/github/lightstep/lightstep-partner-toolkit/tree/main/demo)**.

#### Running locally

* Set the env var `LS_ACCESS_TOKEN` to your Lightstep Access Token
* `yarn install && yarn start`
* Visit `http://localhost:8181`

##### Partner Environment Variables

* `SPLIT_API_TOKEN` - API token for Split.io
* `LD_SDK_KEY` - API token for LaunchDarkly
###  OpenTelemetry + Lightstep Partners Node.js Demo App

Simple Node.js web app that emits OpenTelemetry with Lightstep partner plugins enabled. To preview this in the cloud, **[Run in Codesandbox](https://codesandbox.io/s/github/lightstep/lightstep-partner-toolkit/tree/main/demo)**.

#### Running locally (without Docker)

* Set the env var `LS_ACCESS_TOKEN` to your Lightstep Access Token
* `yarn install && yarn start`
* Visit `http://localhost:8181`

### Running locally (Docker)

* `export LS_ACCESS_TOKEN=<your token>`
* `docker build -t lightstep/lightstep-partner-sdk-demo .`
* `docker run -p 8181:8181 -e LS_ACCESS_TOKEN --rm lightstep/lightstep-partner-sdk-demo`
* Visit `http://localhost:8181`

##### Partner Environment Variables

* `SPLIT_API_TOKEN` - API token for Split.io
* `LD_SDK_KEY` - API token for LaunchDarkly
# new demo env / generatorreceiver - Contributing

## Prerequisites / setup

### Opentelemetry collector builder

Install the [opentelemetry-collector-builder](https://github.com/open-telemetry/opentelemetry-collector-builder); this is deprecated but its replacement does not work with the old version of the collector we're still pinned to.

   1. `$ cd /tmp` (or wherever you like to keep code)
   1. `$ git clone https://github.com/open-telemetry/opentelemetry-collector-builder`
   1. `$ cd opentelemetry-collector-builder`
   1. `$ git checkout v0.35.0`
   1. `$ go get -u golang.org/x/sys`
   1. `$ go install .`

### Get the code

1. Clone [the partner toolkit repo](https://github.com/lightstep/lightstep-partner-toolkit) to a directory of your choosing:

   1.  `$ cd ~/Code` (or wherever)
   1.  `$ git clone https://github.com/lightstep/lightstep-partner-toolkit`
   1.  `$ cd lightstep-partner-toolkit`

1. Check out the development branch - until we break the receiver out into its own repo, this is our effective "main" branch:
    `$ git checkout generatorv2`
1. `$ cd collector` (this will be our working directory for everything that follows)
1. Create a copy of the `hipster_shop.yaml` for local development. Not strictly necessary but will potentially save heartache and hassle 😅 This file is in .gitignore, so it won't be included in your commits. If you want to share config changes, add them to `hipster_shop.yaml` or a new example config file.
   `$ cp generatorreceiver/topos/hipster_shop.yaml generatorreceiver/topos/dev.yaml`

## Environment variables

*best practice would be to use a local env file that works with direnv*

### Access token

You'll need an access token associated with the lightstep project you want to use - your dev project on staging is a good place to start. Go to ⚙ -> Access Tokens to copy an existing one or create a new one. Then:

```shell
$ export LS_ACCESS_TOKEN="<your token>"
```

### Collector endpoint

The env var `OTEL_EXPORTER_OTLP_TRACES_ENDPOINT` determines the endpoint for traces and metrics. If you're using a non-staging project, change this as appropriate.

```shell
$ export OTEL_EXPORTER_OTLP_TRACES_ENDPOINT=ingest.staging.lightstep.com:443
```

### Topo file (generatorreceiver config)

The env var `TOPO_FILE` determines which config file the generatorreceiver uses.

For builder builds you'll want to point to `generatorreceiver/topos/`:

```shell
$ export TOPO_FILE=generatorreceiver/topos/dev.yaml
```

For Docker builds, these files are copied to `/etc/otel/`:

```shell
$ export TOPO_FILE=/etc/otel/dev.yaml
```

## Build and run the collector

There are two options here, but for development purposes I recommend using the builder, which is much quicker and lets you test config changes without rebuilding. With the Docker build method, you need to rebuild the image for all changes, code or config, and the build process takes much longer.

### Build and run with the opentelemetry-collector-builder (recommended)

(You must first install the `opentelemetry-collector-builder`; see Prerequisites above.)

```shell
$ opentelemetry-collector-builder --config ./builder-config.yml
$ /tmp/ls-partner-col-distribution/lightstep-partner-collector --config ./config/collector-config.yml
```

When using the builder, you only need to re-run the first command for code changes; for config changes just re-run the second command. To run with a different topo file, change the `TOPO_FILE` environment variable. This makes this a somewhat faster option for development

If you run into errors while building, ping @Nathan on slack and I can help troubleshoot.

### Build and run with Docker (alternative)

```shell
$ docker build -t lightstep/lightstep-partner-toolkit-collector:latest .
$ docker run --rm -p 8181:8181 -e LS_ACCESS_TOKEN -e OTEL_EXPORTER_OTLP_TRACES_ENDPOINT -e TOPO_FILE lightstep/lightstep-partner-toolkit-collector
```

With Docker, you need to re-run both steps for any code *or* config changes.

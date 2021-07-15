## synthetic-load-generator

Generate OpenTelemetry distributed traces and metrics: no services or code required (just some JSON).

golang port of [Omnition/synthetic-load-generator](https://github.com/Omnition/synthetic-load-generator).

### usage

This creates distributed traces from an architecture defined in a JSON file. See `hipster-shop.json` for an example.

To send data to Lightstep, [sign up for a free account](https://app.lightstep.com/signup) and create an access token by following the instructions [here](https://docs.lightstep.com/docs/create-and-manage-access-tokens).

After a few minutes, you should see metrics and traces in your project.

#### Golang usage

This builds a binary using Go to generate synthetic metrics and traces.

```
    $ go build -o dist/synthetic-load-generator

    # to print traces/metrics to stdout
    $ dist/synthetic-load-generator --paramsFile ./hipster_shop.json

    # to send traces/metrics to a local collector
    $ dist/synthetic-load-generator --paramsFile ./hipster_shop.json --collectorUrl localhost:55680

    # to send traces/metrics to Lightstep
    # get a token for your project [here](https://docs.lightstep.com/docs/create-and-manage-access-tokens).
    $ export LS_ACCESS_TOKEN=<your token>
    $ dist/synthetic-load-generator --paramsFile ./hipster_shop.json --collectorUrl ingest.lightstep.com:443
```

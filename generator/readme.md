## synthentic-load-generator

Generate OpenTelemetry distributed traces: no services or code required (just some JSON).

golang port of [Omnition/synthetic-load-generator](https://github.com/Omnition/synthetic-load-generator).

WIP.

### usage

This creates distributed traces from an architecture defined in a JSON file. See `hipster-shop.json` for an example.

```
    $ go build -o dist/synthetic-load-generator

    # to print traces to stdout
    $ dist/synthetic-load-generator --paramsFile ./hipster_shop.json

    # to send traces to a local collector
    $ dist/synthetic-load-generator --paramsFile ./hipster_shop.json --collectorUrl localhost:55680

    # to send traces to Lightstep
    $ export LS_ACCESS_TOKEN=<your token>
    $ dist/synthetic-load-generator --paramsFile ./hipster_shop.json --collectorUrl ingest.lightstep.com:443
```
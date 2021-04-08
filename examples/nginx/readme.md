## Lightstep + Nginx

Work in progress. Based on [opentelemetry-cpp-contrib](https://github.com/open-telemetry/opentelemetry-cpp-contrib/tree/main/instrumentation/nginx).

#### Starting

```
  $ export LS_ACCESS_TOKEN=<access-token>
  $ docker compose up
  $ curl http://localhost:8000/files/
```
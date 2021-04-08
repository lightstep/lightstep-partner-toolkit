## Lightstep + Nginx

Work in progress. Based on [opentelemetry-cpp-contrib](https://github.com/open-telemetry/opentelemetry-cpp-contrib/tree/main/instrumentation/nginx).

nginx with the OpenTelemetry module enabled + configured to send data to a collector that forwards traces to Lightstep.

#### Running

```
  # Bring up nginx + app server
  $ export LS_ACCESS_TOKEN=<access-token>
  $ docker compose up

  # Make some requests to generate traces!
  $ open http://localhost:8000/
```
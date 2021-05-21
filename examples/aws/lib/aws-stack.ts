import * as cdk from '@aws-cdk/core';
import * as eks from '@aws-cdk/aws-eks';

import { DonutShopApp } from './backends/donut-shop';
import { OtelCollector } from './otel-collector/collector';
import { MockBackend } from './backends/mock-backend';
import { OtelNginxIngress } from './ingress/nginx';
import { Loadtest } from './loadtest/loadtest';
import { LightstepMetricDashboard } from './lightstep/lightstep-metric-dashboard';
import { NginxPlusUpstreamResponseTimeChart } from './lightstep/charts';

export class AwsOtelStack extends cdk.Stack {
  constructor(scope: cdk.Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    if (!process.env.LS_ACCESS_TOKEN) {
      throw 'environment variable LS_ACCESS_TOKEN must be set';
    }

    const cluster = new eks.Cluster(this, 'hello-otel-eks', {
      version: eks.KubernetesVersion.V1_19,
    });

    const donutShop = new DonutShopApp(this, 'donut-shop', {
      cluster,
      serviceName: 'donut-shop-svc',
      servicePort: 80,
      servicePath: '/',
    });

    const coffeeShop = new MockBackend(this, 'coffee', {
      cluster,
      serviceName: 'coffee-svc',
      servicePort: 80,
      servicePath: '/coffee',
    });

    const teaShop = new MockBackend(this, 'tea', {
      cluster,
      serviceName: 'tea-svc',
      servicePort: 80,
      servicePath: '/tea',
    });

    if (process.env.LIGHTSTEP_API_KEY) {
      new LightstepMetricDashboard(this, 'my-dashboard', {
        name: 'NGINX Automatic Dashboard',
        charts: [
          new NginxPlusUpstreamResponseTimeChart(donutShop),
          new NginxPlusUpstreamResponseTimeChart(coffeeShop),
          new NginxPlusUpstreamResponseTimeChart(teaShop),
        ],
        lightstepOrg: 'LightStep',
        lightstepProject: 'Robin-Hipster-Shop',
        lightstepApiKey: process.env.LIGHTSTEP_API_KEY,
      });
    }

    const nginxCollector = new OtelCollector(this, 'otel-collector', {
      cluster,
      serviceName: 'otel-collector-svc',
      servicePort: 80,
      servicePath: '/webhook',
    });
    new Loadtest(this, 'loadtest', {
      cluster,
      targetUrl: 'http://nginx-ingress-svc',
    });

    new OtelNginxIngress(this, 'otel-nginx', {
      cluster: cluster,
      backends: [donutShop, coffeeShop, teaShop, nginxCollector],
    });

    cluster.addHelmChart('Prometheus', {
      chart: 'prometheus',
      repository: 'https://prometheus-community.github.io/helm-charts',
      namespace: 'prometheus',
      values: {
        server: {
          sidecarContainers: [
            {
              name: 'otel-sidecar',
              image: 'lightstep/opentelemetry-prometheus-sidecar',
              imagePullPolicy: 'Always',
              args: [
                '--prometheus.wal=/data/wal',
                '--destination.endpoint=https://ingest.lightstep.com:443',
                '--destination.attribute="service.name=nginx-ingress"',
                `--destination.header=lightstep-access-token=${process.env.LS_ACCESS_TOKEN}`,
              ],
              volumeMounts: [
                {
                  name: 'storage-volume',
                  mountPath: '/data',
                },
              ],
              ports: [
                {
                  name: 'admin-port',
                  containerPort: 9091,
                },
              ],
              livenessProbe: {
                httpGet: {
                  path: '/-/health',
                  port: 'admin-port',
                },
                periodSeconds: 30,
                failureThreshold: 2,
              },
            },
          ],
        },
      },
    });
  }
}

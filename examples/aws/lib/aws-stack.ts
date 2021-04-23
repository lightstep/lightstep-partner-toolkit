import * as cdk from '@aws-cdk/core';
import * as eks from '@aws-cdk/aws-eks';
import { assert } from 'console';
import { readFileSync } from 'fs';

export class AwsStack extends cdk.Stack {
  constructor(scope: cdk.Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);
    
    assert(process.env.LS_ACCESS_TOKEN, 'environment variable LS_ACCESS_TOKEN must be set');

    const cluster = new eks.Cluster(this, 'hello-eks', {
      version: eks.KubernetesVersion.V1_19,
    });

    const appLabel = { app: "donut-shop" };

    const deployment = {
      apiVersion: "apps/v1",
      kind: "Deployment",
      metadata: { name: "donut-shop-pod", namespace: "default" },
      spec: {
        replicas: 1,
        selector: { matchLabels: appLabel },
        template: {
          metadata: { name: "donut-shop", labels: appLabel },
          spec: {
            containers: [
              {
                name: "donut-shop",
                image: "ghcr.io/lightstep/lightstep-partner-toolkit-donut-shop:latest",
                ports: [ { containerPort: 8181 } ],
                env: [
                  { name: 'LS_ACCESS_TOKEN', value: process.env.LS_ACCESS_TOKEN }
                ]
              }
            ]
          }
        }
      }
    };
    
    const service = {
      apiVersion: "v1",
      kind: "Service",
      metadata: { name: "donut-shop-svc", namespace: "default", labels: appLabel },
      spec: {
        type: "ClusterIP",
        ports: [ { name: "http", port: 80, targetPort: 8181 } ],
        selector: appLabel
      }
    };

    const ingress = {
      apiVersion: "networking.k8s.io/v1beta1",
      kind: "Ingress",
      metadata: { 
        name: "hello-kubernetes-ingress",
        namespace: "default"
      },
      spec: {
        ingressClassName: "nginx",
        rules: [
          {
            host: "*.elb.amazonaws.com",
            http: {
              paths: [
                {
                  path: "/",
                  backend: {
                    serviceName: "donut-shop-svc",
                    servicePort: 80
                  }
                }
              ]
            }
          }
        ]
      }
    }

    const configMap = {
      apiVersion: 'v1',
      kind: 'ConfigMap',
      metadata: {
        name: 'otel-config'
      },
      data: {
        "tracer-config.json": readFileSync(`${__dirname}/tracer-config.toml`, 'utf-8')
      }
    }

    cluster.addManifest('hello-kub', service, deployment, ingress, configMap);

    // https://docs.nginx.com/nginx-ingress-controller/overview/
    // Building custom image: https://docs.nginx.com/nginx-ingress-controller/installation/building-ingress-controller-image/
    cluster.addHelmChart("NginxIngress", {
      chart: "nginx-ingress",
      repository: "https://helm.nginx.com/stable",
      namespace: "default",
      values: {
        controller: {
          enableLatencyMetrics: true,
          image: {
             repository: 'smithclay/nginx-ingress-otel',
             tag: '1.11.1',
             pullPolicy: 'Always'
          },
          volumeMounts: [
            {
              name: 'config-opentelemetry',
              mountPath: '/var/lib/nginx/tracer-config.json',
              subPath: 'tracer-config.json'
            }
          ],
          volumes: [
            {
              name: 'config-opentelemetry',
              configMap: {
                name: 'otel-config'
              }
            }
          ]
        },
        prometheus: {
          create: "true",
          port: 9113
        }
      }
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
                `--destination.header=lightstep-access-token=${process.env.LS_ACCESS_TOKEN}`
              ],
              volumeMounts: [
                { 
                  name: 'storage-volume',
                  mountPath: '/data'
                }
              ],
              ports: [
                {
                  name: 'admin-port',
                  containerPort: 9091
                }
              ],
              livenessProbe: {
                httpGet: {
                  path: '/-/health',
                  port: 'admin-port'
                },
                periodSeconds: 30,
                failureThreshold: 2
              }
            }
          ]
        }
      }
    });

  }
}

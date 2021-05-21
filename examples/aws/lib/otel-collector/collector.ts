import * as cdk from '@aws-cdk/core';
import * as eks from '@aws-cdk/aws-eks';
import { readFileSync } from 'fs';
import { NginxBackend, NginxBackendProps } from '../ingress/nginx';

export interface OtelColellectorProps {
  cluster: eks.Cluster;
}

export class OtelCollector extends cdk.Construct implements NginxBackend {
  constructor(
    parent: cdk.Construct,
    name: string,
    public props: NginxBackendProps
  ) {
    super(parent, name);

    const collectorLabel = { app: 'otel-collector' };

    const deployment = {
      apiVersion: 'apps/v1',
      kind: 'Deployment',
      metadata: { name: 'otel-collector', namespace: 'default' },
      spec: {
        replicas: 1,
        selector: { matchLabels: collectorLabel },
        template: {
          metadata: { name: 'otel-collector', labels: collectorLabel },
          spec: {
            containers: [
              {
                name: 'otel-collector',
                image:
                  'ghcr.io/lightstep/lightstep-partner-toolkit-collector:latest',
                imagePullPolicy: 'Always',
                //image: 'otel/opentelemetry-collector-contrib-dev:latest',
                ports: [
                  { containerPort: 55680 },
                  { containerPort: 7070 },
                  { containerPort: 7071 },
                ],
                env: [
                  {
                    name: 'LS_ACCESS_TOKEN',
                    value: process.env.LS_ACCESS_TOKEN,
                  },
                ],
                volumeMounts: [
                  {
                    name: 'otel-collector-config-vol-v3',
                    mountPath: '/etc/otel/config.yaml',
                    subPath: 'config.yaml',
                  },
                ],
              },
            ],
            volumes: [
              {
                name: 'otel-collector-config-vol-v3',
                configMap: {
                  name: 'otel-collector-config-v3',
                },
              },
            ],
          },
        },
      },
    };

    const service = {
      apiVersion: 'v1',
      kind: 'Service',
      metadata: {
        name: 'otel-collector-svc',
        namespace: 'default',
        labels: collectorLabel,
      },
      spec: {
        type: 'ClusterIP',
        ports: [
          { port: 55680, targetPort: 55680, name: 'otel' },
          {
            port: this.props.servicePort,
            targetPort: 7070,
            name: 'http-trace',
          },
          { port: 7071, targetPort: 7071, name: 'http-metric' },
        ],
        selector: collectorLabel,
      },
    };

    const configMap = {
      apiVersion: 'v1',
      kind: 'ConfigMap',
      metadata: {
        name: 'otel-collector-config-v3',
      },
      data: {
        'config.yaml': readFileSync(
          `${__dirname}/partner-collector-config.yaml`,
          //`${__dirname}/collector-config.yaml`,
          'utf-8'
        ),
      },
    };

    props.cluster.addManifest('otel-collector', deployment, service, configMap);
  }
}

import * as cdk from '@aws-cdk/core';
import * as eks from '@aws-cdk/aws-eks';
import { readFileSync } from 'fs';

export interface OtelColellectorProps {
  cluster: eks.Cluster;
}

export class OtelCollector extends cdk.Construct {
  constructor(
    parent: cdk.Construct,
    name: string,
    props: OtelColellectorProps
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
                image: 'otel/opentelemetry-collector-contrib-dev:latest',
                ports: [{ containerPort: 55680 }],
                env: [
                  {
                    name: 'LS_ACCESS_TOKEN',
                    value: process.env.LS_ACCESS_TOKEN,
                  },
                ],
                volumeMounts: [
                  {
                    name: 'otel-collector-config-vol',
                    mountPath: '/etc/otel/config.yaml',
                    subPath: 'config.yaml',
                  },
                ],
              },
            ],
            volumes: [
              {
                name: 'otel-collector-config-vol',
                configMap: {
                  name: 'otel-collector-config',
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
        ports: [{ port: 55680, targetPort: 55680 }],
        selector: collectorLabel,
      },
    };

    const configMap = {
      apiVersion: 'v1',
      kind: 'ConfigMap',
      metadata: {
        name: 'otel-collector-config',
      },
      data: {
        'config.yaml': readFileSync(
          `${__dirname}/collector-config.yaml`,
          'utf-8'
        ),
      },
    };

    props.cluster.addManifest('otel-collector', deployment, service, configMap);
  }
}

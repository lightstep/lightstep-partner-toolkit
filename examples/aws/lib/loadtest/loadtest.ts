import * as cdk from '@aws-cdk/core';
import * as eks from '@aws-cdk/aws-eks';
import { readFileSync } from 'fs';

export interface LoadtestProps {
  cluster: eks.Cluster;
  targetUrl: String;
}

export class Loadtest extends cdk.Construct {
  constructor(
    parent: cdk.Construct,
    name: string,
    props: LoadtestProps
  ) {
    super(parent, name);

    const collectorLabel = { app: 'loadtest' };

    const deployment = {
      apiVersion: 'apps/v1',
      kind: 'Deployment',
      metadata: { name: 'loadtest', namespace: 'default' },
      spec: {
        replicas: 1,
        selector: { matchLabels: collectorLabel },
        template: {
          metadata: { name: 'loadtest', labels: collectorLabel },
          spec: {
            containers: [
              {
                name: 'loadtest-k6',
                image: 'loadimpact/k6:latest',
                volumeMounts: [
                  {
                    name: 'loadtest-config-vol',
                    mountPath: '/tmp/script.js',
                    subPath: 'script.js',
                  },
                ],
                env: [
                  {
                    name: 'TARGET_URL',
                    value: props.targetUrl,
                  },
                ],
                args: [
                  'run', '/tmp/script.js'
                ]
              },
            ],
            volumes: [
              {
                name: 'loadtest-config-vol',
                configMap: {
                  name: 'loadtest-config',
                },
              },
            ],
          },
        },
      },
    };

    const configMap = {
      apiVersion: 'v1',
      kind: 'ConfigMap',
      metadata: {
        name: 'loadtest-config',
      },
      data: {
        'script.js': readFileSync(
          `${__dirname}/script.js`,
          'utf-8'
        ),
      },
    };

    props.cluster.addManifest('loadtest', deployment, configMap);
  }
}

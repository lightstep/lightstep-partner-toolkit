import * as cdk from '@aws-cdk/core';
import * as eks from '@aws-cdk/aws-eks';
import { NginxBackend, NginxBackendProps } from '../ingress/nginx';

export interface DonutShotProps {
  cluster: eks.Cluster;
}

export class DonutShopApp extends cdk.Construct implements NginxBackend {
  constructor(
    parent: cdk.Construct,
    name: string,
    public props: NginxBackendProps
  ) {
    super(parent, name);

    const appLabel = { app: 'donut-shop' };

    const deployment = {
      apiVersion: 'apps/v1',
      kind: 'Deployment',
      metadata: { name: 'donut-shop-pod', namespace: 'default' },
      spec: {
        replicas: 1,
        selector: { matchLabels: appLabel },
        template: {
          metadata: { name: 'donut-shop', labels: appLabel },
          spec: {
            containers: [
              {
                name: 'donut-shop',
                image:
                  'ghcr.io/lightstep/lightstep-partner-toolkit-donut-shop:latest',
                ports: [{ containerPort: 8181 }],
                env: [
                  {
                    name: 'LS_ACCESS_TOKEN',
                    value: process.env.LS_ACCESS_TOKEN,
                  },
                ],
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
        name: this.props.serviceName,
        namespace: 'default',
        labels: appLabel,
      },
      spec: {
        type: 'ClusterIP',
        ports: [
          { name: 'http', port: this.props.servicePort, targetPort: 8181 },
        ],
        selector: appLabel,
      },
    };

    props.cluster.addManifest('donut-shop', deployment, service);
  }
}

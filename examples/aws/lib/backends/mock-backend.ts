import * as cdk from '@aws-cdk/core';
import { NginxBackend, NginxBackendProps } from '../ingress/nginx';

export class MockBackend extends cdk.Construct implements NginxBackend {
  constructor(
    parent: cdk.Construct,
    name: string,
    public props: NginxBackendProps
  ) {
    super(parent, name);

    const deployment = {
      apiVersion: 'apps/v1',
      kind: 'Deployment',
      metadata: { name: this.props.serviceName, namespace: 'default' },
      spec: {
        replicas: 1,
        selector: { matchLabels: { app: this.props.serviceName } },
        template: {
          metadata: { labels: { app: this.props.serviceName } },
          spec: {
            containers: [
              {
                name: this.props.serviceName,
                image: 'nginxdemos/nginx-hello:plain-text',
                ports: [{ containerPort: 8080 }],
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
        labels: { app: this.props.serviceName },
      },
      spec: {
        type: 'ClusterIP',
        ports: [{ port: this.props.servicePort, targetPort: 8080 }],
        selector: { app: this.props.serviceName },
      },
    };

    props.cluster.addManifest(`mock-service-${this.props.serviceName}`, deployment, service);
  }
}

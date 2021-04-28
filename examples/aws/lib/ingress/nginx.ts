import * as cdk from '@aws-cdk/core';
import * as eks from '@aws-cdk/aws-eks';
import { readFileSync } from 'fs';

export interface NginxBackendProps {
  cluster: eks.Cluster;
  serviceName: String;
  servicePort: Number;
  servicePath: String;
}

export interface NginxBackend {
  props: NginxBackendProps;
}

export interface OtelNginxProps {
  cluster: eks.Cluster;
  backends: NginxBackend[];
}

export class OtelNginxIngress extends cdk.Construct {
  constructor(parent: cdk.Construct, name: string, props: OtelNginxProps) {
    super(parent, name);

    const configMap = {
      apiVersion: 'v1',
      kind: 'ConfigMap',
      metadata: {
        name: 'otel-config',
      },
      data: {
        'tracer-config.json': readFileSync(
          `${__dirname}/tracer-config.toml`,
          'utf-8'
        ),
      },
    };

    var imageName = 'smithclay/nginx-ingress-otel';

    if (process.env.DOCKER_CONFIG_BASE64) {
      const secret = {
        apiVersion: 'v1',
        kind: 'Secret',
        metadata: {
          name: 'myregistrykey',
          namespace: 'default'
        },
        data: {
          '.dockerconfigjson': process.env.DOCKER_CONFIG_BASE64
        },
        type: 'kubernetes.io/dockerconfigjson'
      };
      imageName = 'smithclay/nginx-plus-ingress-otel';
      props.cluster.addManifest('otel-nginx-plus', secret)
    }

    props.cluster.addManifest(
      'otel-nginx',
      this.createIngressManifest(props.backends),
      configMap
    );

    // https://docs.nginx.com/nginx-ingress-controller/overview/
    // Building custom image: https://docs.nginx.com/nginx-ingress-controller/installation/building-ingress-controller-image/
    props.cluster.addHelmChart('OTelNginxIngress', {
      chart: 'nginx-ingress',
      repository: 'https://helm.nginx.com/stable',
      namespace: 'default',
      values: {
        controller: {
          nginxplus: process.env.DOCKER_CONFIG_BASE64 ? true : false,
          service: {
            name: 'nginx-ingress-svc'
          },
          serviceAccount: {
            imagePullSecretName: 'myregistrykey'
          },
          enableLatencyMetrics: true,
          image: {
            repository: imageName,
            tag: '1.11.1',
            pullPolicy: 'Always',
          },
          volumeMounts: [
            {
              name: 'config-opentelemetry',
              mountPath: '/var/lib/nginx/tracer-config.json',
              subPath: 'tracer-config.json',
            },
          ],
          volumes: [
            {
              name: 'config-opentelemetry',
              configMap: {
                name: 'otel-config',
              },
            },
          ],
        },
        prometheus: {
          create: 'true',
          port: 9113,
        },
      },
    });
  }

  createIngressManifest(backends: NginxBackend[]): any {
    const ingress: any = {
      apiVersion: 'networking.k8s.io/v1beta1',
      kind: 'Ingress',
      metadata: {
        name: 'nginx-ingress',
        namespace: 'default',
      },
      spec: {
        ingressClassName: 'nginx',
        rules: [
          {
            host: '*.elb.amazonaws.com',
            http: {
              paths: [],
            },
          },
        ],
      },
    };

    for (const b of backends) {
      ingress.spec.rules[0].http.paths.push({
        path: b.props.servicePath,
        backend: {
          serviceName: b.props.serviceName,
          servicePort: b.props.servicePort,
        },
      });
    }

    return ingress;
  }
}

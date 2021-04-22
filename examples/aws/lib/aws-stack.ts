import * as cdk from '@aws-cdk/core';
import * as eks from '@aws-cdk/aws-eks';

export class AwsStack extends cdk.Stack {
  constructor(scope: cdk.Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    const cluster = new eks.Cluster(this, 'hello-eks', {
      version: eks.KubernetesVersion.V1_19,
    });

    const appLabel = { app: "hello-kubernetes" };

    const deployment = {
      apiVersion: "apps/v1",
      kind: "Deployment",
      metadata: { name: "hello-kubernetes-pod", namespace: "default" },
      spec: {
        replicas: 1,
        selector: { matchLabels: appLabel },
        template: {
          metadata: { name: "hello-kubernetes", labels: appLabel },
          spec: {
            containers: [
              {
                name: "hello-kubernetes",
                image: "paulbouwer/hello-kubernetes:1.5",
                ports: [ { containerPort: 8080 } ]
              }
            ]
          }
        }
      }
    };
    
    const service = {
      apiVersion: "v1",
      kind: "Service",
      metadata: { name: "hello-kubernetes-svc", namespace: "default", labels: appLabel },
      spec: {
        type: "ClusterIP",
        ports: [ { name: "http", port: 80, targetPort: 8080 } ],
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
                    serviceName: "hello-kubernetes-svc",
                    servicePort: 80
                  }
                }
              ]
            }
          }
        ]
      }
    }

    cluster.addManifest('hello-kub', service, deployment, ingress);

    // https://docs.nginx.com/nginx-ingress-controller/overview/
    // Building custom image: https://docs.nginx.com/nginx-ingress-controller/installation/building-ingress-controller-image/
    cluster.addHelmChart("NginxIngress", {
      chart: "nginx-ingress",
      repository: "https://helm.nginx.com/stable",
      namespace: "default",
      values: {
        controller: {
          enableLatencyMetrics: true
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

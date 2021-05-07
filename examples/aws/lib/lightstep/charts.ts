import { NginxBackend } from '../ingress/nginx';

export interface LightstepMetricChart {
  toJSON(rank: number): object;
}

export class NginxPlusUpstreamResponseTimeChart
  implements LightstepMetricChart {
  constructor(protected backend: NginxBackend, protected title?: string) {}

  toJSON(rank: number) {
    return {
      rank: rank,
      title:
        this.title ||
        `Upstream Server Response Time (ms): ${this.backend.props.serviceName}`,
      'chart-type': 'timeseries',
      'metric-queries': [
        {
          'query-name': 'a',
          'query-type': 'single',
          hidden: false,
          'display-type': 'line',
          'metric-query': {
            metric: 'nginx_ingress_nginxplus_upstream_server_response_time',
            filters: [
              {
                key: 'upstream',
                value: `default-nginx-ingress-*.elb.amazonaws.com-${this.backend.props.serviceName}-${this.backend.props.servicePort}`,
                operand: 'eq',
              },
            ],
            'timeseries-operator': 'last',
            'group-by': {
              'label-keys': [],
              'aggregation-method': 'avg',
            },
          },
        },
      ],
    };
  }
}

import lambda = require('@aws-cdk/aws-lambda');
import { Construct, CustomResource, Duration } from '@aws-cdk/core';
import { RetentionDays } from '@aws-cdk/aws-logs';
import { Provider } from '@aws-cdk/custom-resources';

import fs = require('fs');

export interface LightstepDashboardProps {
  name: string;
  lightstepProject: string;
  lightstepOrg: string;
  lightstepApiKey: string;
}

export class LightstepMetricDashboard extends Construct {
  public readonly response: string;

  constructor(scope: Construct, id: string, props: LightstepDashboardProps) {
    super(scope, id);

    const onEvent = new lambda.SingletonFunction(this, 'Singleton', {
      uuid: 'f7d4f730-4ee1-11e8-8c2d-da7ae01bbebc',
      code: new lambda.InlineCode(
        fs.readFileSync(`${__dirname}/custom-resource-handler.py`, {
          encoding: 'utf-8',
        })
      ),
      environment: {
        LIGHTSTEP_API_KEY: props.lightstepApiKey,
      },
      handler: 'index.main',
      timeout: Duration.seconds(300),
      runtime: lambda.Runtime.PYTHON_3_6,
    });

    const myProvider = new Provider(this, 'LightstepMetricDashboardProvider', {
      onEventHandler: onEvent,
      logRetention: RetentionDays.ONE_DAY,
    });
    const resource = new CustomResource(this, `LightstepDash-${id}`, {
      serviceToken: myProvider.serviceToken,
      properties: {
        name: props.name,
        lightstepOrg: props.lightstepOrg,
        lightstepProject: props.lightstepProject,
      },
    });

    this.response = resource.getAtt('Response').toString();
  }
}

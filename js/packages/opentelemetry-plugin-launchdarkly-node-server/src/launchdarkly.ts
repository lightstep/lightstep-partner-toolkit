import { BasePlugin } from '@opentelemetry/core';
import { SpanKind } from '@opentelemetry/api';
import type * as ldTypes from 'launchdarkly-node-server-sdk';

import * as shimmer from 'shimmer';
import { VERSION } from './version';

export class LaunchDarklyNodeServerPlugin extends BasePlugin<typeof ldTypes> {
  static readonly COMPONENT = 'launchdarkly-node-server-sdk';
  static readonly COMPONENT_SHORT = 'launchdarkly-node';

  readonly supportedVersions = ['5.14.1', '5.14.1'];

  constructor(readonly moduleName: string) {
    super('@opentelemetry/plugin-launchdarkly', VERSION);
  }

  protected patch() {
    if (this._moduleExports) {
      this._logger.debug('patching launchdarkly');
      shimmer.wrap(this._moduleExports, 'init', this._getInitPatch.bind(this));
    }
    return this._moduleExports;
  }

  protected unpatch(): void {
    if (this._moduleExports) {
      shimmer.unwrap(this._moduleExports, 'init');
    }
  }

  private _getVariationPatch(original: Function) {
    const instrumentation = this;
    instrumentation._logger.debug('Patching variation function');
    return async function variation(this: any) {
      const span = instrumentation._tracer.startSpan(
        'launchdarkly - variation',
        {
          kind: SpanKind.CLIENT,
          attributes: {
            component: LaunchDarklyNodeServerPlugin.COMPONENT_SHORT,
          },
        }
      );
      const value = await original.apply(this, arguments);
      span.setAttributes({
        'launchdarkly.com.flag': arguments[0],
        'launchdarkly.com.value': value,
        'launchdarkly.com.user': JSON.stringify(arguments[1]),
      });
      span.end();
      return value;
    };
  }

  private _getInitPatch(original: (options?: any) => any) {
    const instrumentation = this;
    return function init(this: any, opts?: any) {
      const newClient = original.apply(this, [opts]);

      shimmer.wrap(
        newClient,
        'variation',
        instrumentation._getVariationPatch.bind(instrumentation)
      );

      return newClient;
    };
  }
}

export const plugin = new LaunchDarklyNodeServerPlugin(
  LaunchDarklyNodeServerPlugin.COMPONENT
);

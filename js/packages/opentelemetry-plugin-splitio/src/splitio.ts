import { BasePlugin } from '@opentelemetry/core';
import { SpanKind } from '@opentelemetry/api';
import type * as splitioTypes from '@splitsoftware/splitio';

import * as shimmer from 'shimmer';
import { VERSION } from './version';

export class SplitioPlugin extends BasePlugin<typeof splitioTypes> {
  static readonly COMPONENT = '@splitsoftware/splitio';
  static readonly COMPONENT_SHORT = 'splitio';

  readonly supportedVersions = ['10.15.2', '10.15.2'];

  constructor(readonly moduleName: string) {
    super('@opentelemetry/plugin-splitio', VERSION);
  }

  protected patch() {
    if (this._moduleExports) {
      this._logger.debug('patching splitio');
      shimmer.wrap(
        this._moduleExports,
        'SplitFactory',
        this._getFactoryPatch.bind(this)
      );
    }
    return this._moduleExports;
  }

  protected unpatch(): void {
    if (this._moduleExports) {
      shimmer.unwrap(this._moduleExports.SplitFactory.prototype, 'client');
    }
  }

  private _getTreatmentPatch(original: Function) {
    const instrumentation = this;
    instrumentation._logger.debug('Patching getTreatment function');
    return async function getTreatment(this: any) {
      const span = instrumentation._tracer.startSpan('splitio - getTreatment', {
        kind: SpanKind.CLIENT,
        attributes: { component: SplitioPlugin.COMPONENT_SHORT },
      });
      const treatment = await original.apply(this, arguments);
      span.setAttributes({
        'split.io.key': arguments[0],
        'split.io.value': treatment,
        'split.io.treatment': arguments[1],
      });
      span.end();
      return treatment;
    };
  }

  private _getFactoryClientPatch(original: (options?: any) => any) {
    const instrumentation = this;

    return function client(this: any, opts?: any) {
      const newClient = original.apply(this, [opts]);
      shimmer.wrap(
        newClient,
        'getTreatment',
        instrumentation._getTreatmentPatch.bind(instrumentation)
      );

      return newClient;
    };
  }

  private _getFactoryPatch(original: (options?: any) => any) {
    const instrumentation = this;
    return function factory(this: any, opts?: any) {
      const newFactory = original.apply(this, [opts]);

      shimmer.wrap(
        newFactory,
        'client',
        instrumentation._getFactoryClientPatch.bind(instrumentation)
      );

      return newFactory;
    };
  }
}

export const plugin = new SplitioPlugin(SplitioPlugin.COMPONENT);

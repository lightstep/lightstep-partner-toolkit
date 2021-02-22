import { BasePlugin } from '@opentelemetry/core';
import { getSpan, context } from '@opentelemetry/api';
import * as Rollbar from 'rollbar';

import * as shimmer from 'shimmer';
import { VERSION } from './version';

export class RollbarPlugin extends BasePlugin<typeof Rollbar> {
  static readonly COMPONENT = 'rollbar';
  static readonly COMPONENT_SHORT = 'rollbar';

  readonly supportedVersions = ['2.19.4', '2.19.4'];

  constructor(readonly moduleName: string) {
    super('rollbar', VERSION);
  }

  protected patch() {
    if (this._moduleExports) {
      this._logger.debug('patching rollbar');
      shimmer.wrap(
        this._moduleExports.prototype,
        'error',
        this._getErrorPatch.bind(this)
      );
    }
    return this._moduleExports;
  }

  protected unpatch(): void {
    if (this._moduleExports) {
      shimmer.unwrap(this._moduleExports.prototype, 'error');
    }
  }

  private _getErrorPatch(original: (...args: Rollbar.LogArgument[]) => Rollbar.LogResult) {
    return function error(this: any, ...args: Rollbar.LogArgument[]) : Rollbar.LogResult {
      const span = getSpan(context.active());
      const result = original.apply(this, args);
      if (span) {
        span.setAttributes({
          'rollbar.has_error': true,
          'rollbar.uuid': result.uuid,
          'error': true
        });
      }
      return result
    }
  }    
}

export const plugin = new RollbarPlugin(RollbarPlugin.COMPONENT);

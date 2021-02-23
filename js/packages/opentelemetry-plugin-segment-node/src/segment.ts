import { BasePlugin } from '@opentelemetry/core';
import { getSpan, context, diag } from '@opentelemetry/api';
const Analytics = require('analytics-node');

import * as shimmer from 'shimmer';
import { VERSION } from './version';

export class SegmentPlugin extends BasePlugin<typeof Analytics> {
  static readonly COMPONENT = 'analytics-node';
  static readonly COMPONENT_SHORT = 'analytics-node';

  readonly supportedVersions = ['4.0.0', '4.0.0'];

  constructor(readonly moduleName: string) {
    super('analytics-node', VERSION);
  }

  protected patch() {
    if (this._moduleExports) {
      diag.debug('patching analytics-node');
      shimmer.wrap(
        this._moduleExports.prototype,
        'track',
        this._getTrackPatch.bind(this)
      );
    }
    return this._moduleExports;
  }

  protected unpatch(): void {
    if (this._moduleExports) {
      shimmer.unwrap(this._moduleExports.prototype, 'track');
    }
  }

  private _getTrackPatch(original: Function) {
    return function track(this: any, message?: any, callback?: any): any {
      const span = getSpan(context.active());
      const result = original.apply(this, [message, callback]);
      if (span && message.event) {
        span.addEvent(`segment.io track - ${message.event}`, {
          userId: message.userId,
          anonymousId: message.anonymousId,
          ...message.properties,
        });
      }
      return result;
    };
  }
}

export const plugin = new SegmentPlugin(SegmentPlugin.COMPONENT);

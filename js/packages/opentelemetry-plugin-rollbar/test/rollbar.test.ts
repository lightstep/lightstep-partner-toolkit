/*
 * Copyright The OpenTelemetry Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

import {
  context,
  NoopLogger,
  setSpan,
} from '@opentelemetry/api';
import { NodeTracerProvider } from '@opentelemetry/node';
import { AsyncHooksContextManager } from '@opentelemetry/context-async-hooks';
import {
  InMemorySpanExporter,
  SimpleSpanProcessor,
} from '@opentelemetry/tracing';
import * as assert from 'assert';
import * as Rollbar from 'rollbar';
import { plugin, RollbarPlugin } from '../src';

const memoryExporter = new InMemorySpanExporter();

describe('rollbar@2.19.x', () => {
  const provider = new NodeTracerProvider();
  const tracer = provider.getTracer('external');

  let contextManager: AsyncHooksContextManager;
  let rollbar: typeof Rollbar;
  beforeEach(() => {
    contextManager = new AsyncHooksContextManager().enable();
    context.setGlobalContextManager(contextManager);
  });

  afterEach(() => {
    context.disable();
  });

  before(function () {
    rollbar = require('rollbar');
    provider.addSpanProcessor(new SimpleSpanProcessor(memoryExporter));
    plugin.enable(rollbar, provider, new NoopLogger());
  });

  it('should have correct module name', () => {
    assert.strictEqual(plugin.moduleName, RollbarPlugin.COMPONENT);
  });

  describe('#error()', () => {
    it('should set attributes on calls to error', done => {
      let rollbarClient = new rollbar({})

      const span = tracer.startSpan('test span');
      context.with(setSpan(context.active(), span), () => {
        const span = tracer.startSpan('error span');
        context.with(setSpan(context.active(), span), () => {
          rollbarClient.error('this is an error');
        })
        span.end();
      });
      const endedSpans = memoryExporter.getFinishedSpans();
      assert.strictEqual(endedSpans.length, 1);
      assert.strictEqual(
        endedSpans[0].attributes['rollbar.has_error'],
        true
      );
      assert.strictEqual(
        endedSpans[0].attributes['error'],
        true
      );
      done();
    });
  });

  describe('Removing instrumentation', () => {
    before(() => {
      memoryExporter.reset()
      plugin.disable();
    });

    it(`should not create a child span`, done => {
      let rollbarClient = new rollbar({})
      const span = tracer.startSpan('test span');
      context.with(setSpan(context.active(), span), () => {
        const span = tracer.startSpan('error span');
        context.with(setSpan(context.active(), span), () => {
          rollbarClient.error('this is an error');
        })
        span.end();
      });
      const endedSpans = memoryExporter.getFinishedSpans();
      assert.strictEqual(endedSpans.length, 1);
      assert.strictEqual(
        endedSpans[0].attributes['rollbar.has_error'],
        undefined
      );
      done();
    });
  });
});
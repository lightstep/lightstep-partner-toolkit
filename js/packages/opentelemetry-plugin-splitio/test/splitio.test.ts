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
import type * as splitioTypes from '@splitsoftware/splitio';
import { plugin, SplitioPlugin } from '../src';

const memoryExporter = new InMemorySpanExporter();


describe('@splitsoftware/splitio@10.x', () => {
  const provider = new NodeTracerProvider();
  const tracer = provider.getTracer('external');

  let contextManager: AsyncHooksContextManager;
  let splitio: typeof splitioTypes;
  beforeEach(() => {
    contextManager = new AsyncHooksContextManager().enable();
    context.setGlobalContextManager(contextManager);
  });

  afterEach(() => {
    context.disable();
  });

  before(function () {
    splitio = require('@splitsoftware/splitio');
    provider.addSpanProcessor(new SimpleSpanProcessor(memoryExporter));
    plugin.enable(splitio, provider, new NoopLogger());
  });

  it('should have correct module name', () => {
    assert.strictEqual(plugin.moduleName, SplitioPlugin.COMPONENT);
  });

  describe('#getTreatment()', () => {
    it('should propagate the current span to getTreatment', done => {
      const span = tracer.startSpan('test span');
      let splitSdk: SplitIO.ISDK;
      let splitClient: SplitIO.IClient;

      context.with(setSpan(context.active(), span), () => {
        splitSdk = splitio.SplitFactory({ core: { authorizationKey: 'TEST_KEY' } });
        splitClient = splitSdk.client();
        const span = tracer.startSpan('test span');
        context.with(setSpan(context.active(), span), async () => {
          const treatment = await splitClient.getTreatment('customer-1234', 'TEST_SPLIT');
          const endedSpans = memoryExporter.getFinishedSpans();

          assert.strictEqual(endedSpans.length, 1);
          assert.strictEqual(
            endedSpans[0].attributes['split.io.key'],
            'customer-1234'
          );
          assert.strictEqual(
            endedSpans[0].attributes['split.io.treatment'],
            'TEST_SPLIT'
          );
          assert.strictEqual(
            endedSpans[0].attributes['split.io.value'],
            treatment
          );
          assert.strictEqual(
            endedSpans[0].name,
            `splitio - getTreatment`
          );
          done();
        })
      });
    });
  });
});
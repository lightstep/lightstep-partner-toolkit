import { context, setSpan } from '@opentelemetry/api';
import { NodeTracerProvider } from '@opentelemetry/node';
import { AsyncHooksContextManager } from '@opentelemetry/context-async-hooks';
import {
  InMemorySpanExporter,
  SimpleSpanProcessor,
} from '@opentelemetry/tracing';
import * as assert from 'assert';
const Analytics = require('analytics-node');

import { plugin, SegmentPlugin } from '../src';

const memoryExporter = new InMemorySpanExporter();

describe('segment@2.19.x', () => {
  const provider = new NodeTracerProvider();
  const tracer = provider.getTracer('external');

  let contextManager: AsyncHooksContextManager;
  let analytics: typeof Analytics;
  beforeEach(() => {
    contextManager = new AsyncHooksContextManager().enable();
    context.setGlobalContextManager(contextManager);
  });

  afterEach(() => {
    context.disable();
  });

  before(() => {
    analytics = require('analytics-node');
    provider.addSpanProcessor(new SimpleSpanProcessor(memoryExporter));
    plugin.enable(analytics, provider);
  });

  it('should have correct module name', () => {
    assert.strictEqual(plugin.moduleName, SegmentPlugin.COMPONENT);
  });

  describe('#track()', () => {
    it('should add events on calls to track', done => {
      const segmentClient = new analytics('write key');

      const span = tracer.startSpan('test span');
      context.with(setSpan(context.active(), span), () => {
        const span = tracer.startSpan('user interaction span');
        context.with(setSpan(context.active(), span), () => {
          segmentClient.track({ event: 'button clicked', userId: 'foo_123' });
        });
        span.end();
      });
      const endedSpans = memoryExporter.getFinishedSpans();
      assert.strictEqual(endedSpans.length, 1);
      assert.strictEqual(endedSpans[0].events.length, 1);
      assert.strictEqual(
        endedSpans[0].events[0].name,
        'segment.io track - button clicked'
      );
      done();
    });
  });

  describe('Removing instrumentation', () => {
    before(() => {
      memoryExporter.reset();
      plugin.disable();
    });

    it('should not create a child span', done => {
      const segmentClient = new analytics({});
      const span = tracer.startSpan('test span');
      context.with(setSpan(context.active(), span), () => {
        const span = tracer.startSpan('error span');
        context.with(setSpan(context.active(), span), () => {
          segmentClient.track({ event: 'button clicked', userId: 'foo_123' });
        });
        span.end();
      });
      const endedSpans = memoryExporter.getFinishedSpans();
      assert.strictEqual(endedSpans.length, 1);
      assert.strictEqual(endedSpans[0].events.length, 0);
      done();
    });
  });
});

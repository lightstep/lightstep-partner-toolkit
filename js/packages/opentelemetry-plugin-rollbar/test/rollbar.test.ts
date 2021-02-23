import { context, setSpan } from '@opentelemetry/api';
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

  before(() => {
    rollbar = require('rollbar');
    const config = {
      // TODO: add plugin options here once supported
    };
    provider.addSpanProcessor(new SimpleSpanProcessor(memoryExporter));
    plugin.enable(rollbar, provider, config);
  });

  it('should have correct module name', () => {
    assert.strictEqual(plugin.moduleName, RollbarPlugin.COMPONENT);
  });

  describe('#error()', () => {
    it('should set attributes on calls to error', done => {
      const rollbarClient = new rollbar({});

      const span = tracer.startSpan('test span');
      context.with(setSpan(context.active(), span), () => {
        const span = tracer.startSpan('error span');
        context.with(setSpan(context.active(), span), () => {
          rollbarClient.error('this is an error');
        });
        span.end();
      });
      const endedSpans = memoryExporter.getFinishedSpans();
      assert.strictEqual(endedSpans.length, 1);
      assert.strictEqual(endedSpans[0].attributes['rollbar.has_error'], true);
      assert.strictEqual(endedSpans[0].attributes['error'], true);
      done();
    });
  });

  describe('Removing instrumentation', () => {
    before(() => {
      memoryExporter.reset();
      plugin.disable();
    });

    it('should not create a child span', done => {
      const rollbarClient = new rollbar({});
      const span = tracer.startSpan('test span');
      context.with(setSpan(context.active(), span), () => {
        const span = tracer.startSpan('error span');
        context.with(setSpan(context.active(), span), () => {
          rollbarClient.error('this is an error');
        });
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

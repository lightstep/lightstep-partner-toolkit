const { SplitFactory } = require('@splitsoftware/splitio');
const LaunchDarkly = require('launchdarkly-node-server-sdk');

let splitClient;
let ldClient;
/**
 * Split.io Node SDK Initialization
 */
if (process.env.SPLIT_API_KEY) {
  const factory = SplitFactory({
    core: {
      authorizationKey: process.env.SPLIT_API_KEY,
    },
  });

  const sc = factory.client();
  sc.on(sc.Event.SDK_READY, () => {
    // eslint-disable-next-line no-console
    console.log('split.io sdk is ready');
  });
  splitClient = sc;
}
/**
 * LaunchDarkly Node (Server) SDK Initialization
 */
if (process.env.LD_SDK_KEY) {
  const ld = LaunchDarkly.init(process.env.LD_SDK_KEY);
  ld.once('ready', () => {
    // eslint-disable-next-line no-console
    console.log('launchdarkly sdk is ready');
  });
  ldClient = ld;
}

/**
 * Wraps multiple feature flag SDKs
 */
module.exports.getFeatureFlag = async (customerId, flagName) => {
  if (ldClient) {
    return ldClient.variation(flagName, { key: customerId }, false);
  }

  if (splitClient) {
    return splitClient.getTreatment(
      customerId,
      flagName,
    );
  }

  return null;
};

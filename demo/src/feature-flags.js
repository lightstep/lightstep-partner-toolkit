const SplitFactory = require("@splitsoftware/splitio").SplitFactory;
const LaunchDarkly = require('launchdarkly-node-server-sdk');

var splitClient;
var ldClient;
/**
 * Split.io Node SDK Initialization
 */
if (process.env.SPLIT_API_KEY) {
  var factory = SplitFactory({
    core: {
      authorizationKey: process.env.SPLIT_API_KEY
    }
  });
  
  const sc = factory.client();
  sc.on(splitClient.Event.SDK_READY, function () {
    console.log("split.io sdk is ready");
  });
  splitClient = sc;
}
/**
 * LaunchDarkly Node (Server) SDK Initialization
 */
if (process.env.LD_SDK_KEY) {
  const ld = LaunchDarkly.init(process.env.LD_SDK_KEY);
  ld.once("ready", () => {
    console.log("launchdarkly sdk is ready")
  });
  ldClient = ld;
}

/**
 * Wraps multiple feature flag SDKs
 */
module.exports.getFeatureFlag = async (customerId, flagName) => {
  if (ldClient) {
    return await ldClient.variation(flagName, { key : customerId }, false); 
  }

  if (splitClient) {
    return await splitClient.getTreatment(
      customerId,
      flagName
    );
  }

  return null;
};
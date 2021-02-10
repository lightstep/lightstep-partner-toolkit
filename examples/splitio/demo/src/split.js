var SplitFactory = require("@splitsoftware/splitio").SplitFactory;

var factory = SplitFactory({
  core: {
    authorizationKey: process.env.SPLIT_API_KEY
  }
});

const splitClient = factory.client();
splitClient.on(splitClient.Event.SDK_READY, function () {
  console.log("split.io sdk is ready");
});

module.exports = splitClient;

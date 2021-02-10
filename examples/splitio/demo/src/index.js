// Important! Require this code before everything else
// This configures tracing of all following dependencies
const serviceName = process.env.LS_SERVICE_NAME || "splitio-test-app";
require("./tracer")(serviceName);

let express = require("express");
let splitClient = require("./split");

var cors = require('cors')
let mustacheExpress = require('mustache-express');

let app = express();
const PORT = process.env.PORT || 8181;

// TODO: allow setting tracer._tracerProvider._registeredSpanProcessors[0]._exporter.headers["Lightstep-Access-Token"]
// via web interface.

const router = express.Router();
router.get("/donuts", async (req, res) => {
  let result = await splitClient.getTreatment(
    Math.floor(Math.random() * 1000), // random user
    "DONUT_EXPERIMENT"
  );
  console.log("split result", result);
  if (result === "on") {
    // simulate slow response
    setTimeout(() => {
      return res.json({ donuts: [{ type: "sprinkles" }] });
    }, 1000);
  } else {
    // simulate fast response
    setTimeout(() => {
      return res.json({ donuts: [{ type: "chocolate" }] });
    }, 300);
  }
});

app.set('views', `${__dirname}/views`);
app.engine('mustache', mustacheExpress());
app.set('view engine', 'mustache');

app.use(cors())
app.use("/api", router);

app.get('/', function(req, res) {
  const data = { LS_SERVICE_NAME: `${serviceName}-browser` };
  if (process.env.LS_ACCESS_TOKEN) {
    data.LS_ACCESS_TOKEN = process.env.LS_ACCESS_TOKEN
  }

  res.render('index', data);
});

app.listen(PORT, () => {
  console.log(`server is live and listening to ${PORT}`);
});

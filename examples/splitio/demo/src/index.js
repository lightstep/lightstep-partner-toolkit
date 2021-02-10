require("./tracer")("splitio-test-app");

let express = require("express");
let splitClient = require("./split");

let app = express();
const PORT = 8181;

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

app.use("/api", router);
app.use(express.static("public"));

app.listen(PORT, () => {
  console.log(`server is live and listening to ${PORT}`);
});

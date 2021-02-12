// Important! Require this code before everything else
// This configures tracing of all following dependencies
const serviceName = process.env.LS_SERVICE_NAME || 'lightstep-partner-test-app';
require('./tracer')(serviceName);

const express = require('express');

const cors = require('cors');
const mustacheExpress = require('mustache-express');

const app = express();
const apiRouter = require('./router');

const PORT = process.env.PORT || 8181;

app.engine('mustache', mustacheExpress());
app.set('views', `${__dirname}/views`);
app.set('view engine', 'mustache');

app.use(cors());
app.use('/api', apiRouter);

app.get('/', (req, res) => {
  const data = {
    LS_SERVICE_NAME: `${serviceName}-browser`,
    LS_PROJECT_NAME: process.env.LS_PROJECT_NAME,
    LS_ACCESS_TOKEN: process.env.LS_ACCESS_TOKEN,
  };

  res.render('index', data);
});

app.listen(PORT, () => {
  // eslint-disable-next-line no-console
  console.log(`server is live and listening to ${PORT}`);
});

const express = require('express');
const Rollbar = require('rollbar');

const rollbar = new Rollbar({ accessToken: process.env.ROLLBAR_POST_ITEM_KEY });

const router = express.Router();
const { getFeatureFlag } = require('./feature-flags');

/**
 * Returns a type of donut based on a feature flag
 * for a given customer.
 */
router.get('/donuts', async (req, res) => {
  const customerId = Math.floor(Math.random() * 1000);
  const result = await getFeatureFlag(
    customerId,
    'DONUT_EXPERIMENT',
  );

  // Error ~25% of the time
  if (Math.floor(Math.random() * 4) === 3) {
    rollbar.error('Undercooked Donut Error');
  }

  if (result === 'on') {
    // simulate slow response
    setTimeout(() => res.json({ donuts: [{ type: 'sprinkles' }] }), 1000);
  } else {
    // simulate fast response
    setTimeout(() => res.json({ donuts: [{ type: 'chocolate' }] }), 300);
  }
});

module.exports = router;

const express = require('express');

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

  if (result === 'on') {
    // simulate slow response
    setTimeout(() => res.json({ donuts: [{ type: 'sprinkles' }] }), 1000);
  } else {
    // simulate fast response
    setTimeout(() => res.json({ donuts: [{ type: 'chocolate' }] }), 300);
  }
});

module.exports = router;

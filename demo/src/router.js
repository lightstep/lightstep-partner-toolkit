const express = require('express');
const Rollbar = require('rollbar');
const Analytics = require('analytics-node');
const AWS = require('aws-sdk');

const s3 = new AWS.S3({ apiVersion: '2006-03-01' });

const rollbar = new Rollbar({ accessToken: process.env.ROLLBAR_POST_ITEM_KEY });
const analytics = new Analytics(process.env.SEGMENT_WRITE_KEY || 'write-key', { enable: false });

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

  analytics.track({
    userId: customerId || 'unknown',
    event: 'Donut Ordered',
    properties: {
      donut_type: result,
    },
  });

  // Error ~25% of the time
  if (Math.floor(Math.random() * 4) === 3) {
    rollbar.error('Undercooked Donut Error');
    res.status(500);
    return res.json({ error: 'undercooked donuts' });
  }

  if (result === 'on') {
    // simulate slow response: calls out to public S3 bucket
    s3.listObjectsV2({ Bucket: 'njogis-imagery' }, (err, data) => {
      if (err) {
        return rollbar.error(err);
      }
      return res.json({ donuts: [{ type: 'sprinkles' }], s3: data });
    });
  } else {
    // simulate fast response
    setTimeout(() => res.json({ donuts: [{ type: 'chocolate' }] }), 300);
  }
});

module.exports = router;

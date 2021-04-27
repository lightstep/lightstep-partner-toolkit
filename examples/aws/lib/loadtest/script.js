import http from "k6/http";

// Options
export let options = {
    stages: [
        // Linearly ramp up from 1 to 50 VUs during first minute
        { target: 50, duration: "1m" },
        // Hold at 50 VUs for the next 3 minutes and 30 seconds
        { target: 50, duration: "3m30s" },
        // Linearly ramp down from 50 to 0 50 VUs over the last 30 seconds
        { target: 0, duration: "30s" }
        // Total execution time will be ~5 minutes
    ]
};

export default function() {
    http.get(__ENV.TARGET_URL);
    http.get(`${__ENV.TARGET_URL}/coffee`);
    http.get(`${__ENV.TARGET_URL}/tea`);
    http.get(`${__ENV.TARGET_URL}/api/donuts`);
};

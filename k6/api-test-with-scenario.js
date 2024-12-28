// import necessary module
import { check } from 'k6';
import http from 'k6/http';

export const options = {
  thresholds: {
    http_req_failed: ['rate<0.01'], // http errors should be less than 1%
    http_req_duration: ['p(99)<1000'] // 99% of requests should be below 1000ms = 1s
  },
  // define scenarios
  // if added scenario we do not need to pass --iteration opts
  scenarios: {
    // arbitrary name of scenario
    average_load: {
      executor: 'ramping-vus',
      stages: [
        // ramp up to average load of 20 virtual users
        { duration: '10s', target: 20 },
        // maintain load
        { duration: '50s', target: 20 },
        // ramp down to zero
        { duration: '5s', target: 0 },
      ],
    },
  },
};

export default function () {
  const randomPath = Math.floor(Math.random() * 1000) + 1;
  const randomUrl = `https://example.com/${randomPath}`;

  const url = 'http://localhost:3002/';
  const payload = JSON.stringify({
   url: randomUrl,
  });

  const params = {
    headers: {
      'Content-Type': 'application/json',
    },
  };

  // send a post request and save response as a variable
  const res = http.post(url, payload, params);
  check(res, {
    'response code was 200': (res) => res.status == 200,
  })
}

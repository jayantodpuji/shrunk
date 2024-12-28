## Introduction

The url-shortener api here is just a test media for me to learn K6 and how i can improve my url-shortener development.

## Test#1

Scenario use:

```js
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
```

In this scenario, i got error for one of my threshold which is the http_req_failed.
```
 http_req_failed................: 97.23% 34846 out of 35836
```

This is because the random URL generator in the test produced same url, so the backend generate same slug.

```go
func shrunk(originalUrl string) string {
	hash := sha1.Sum([]byte(originalUrl))
	return base64.RawURLEncoding.EncodeToString(hash[:])[:7]
}
```

The expected behavior is if the request url is already exist in database, the backend should generate different slug.

### Iteration#1

Add counter if the error is violated the primary key constraint since the slug is a `pkey`. You can see the changes at this [commit](https://github.com/jayantodpuji/shrunk/commit/8434e20c581be037e73689a38cb93582da1e3954)

Result: It pass the threshold.
```
     ✓ response code was 200

     checks.........................: 100.00% 11929 out of 11929
     data_received..................: 1.5 MB  23 kB/s
     data_sent......................: 2.0 MB  30 kB/s
     http_req_blocked...............: avg=7.17µs  min=2µs    med=5µs     max=1.44ms p(90)=8µs      p(95)=10µs
     http_req_connecting............: avg=538ns   min=0s     med=0s      max=642µs  p(90)=0s       p(95)=0s
   ✓ http_req_duration..............: avg=97.21ms min=1.88ms med=58.76ms max=2.21s  p(90)=215.72ms p(95)=272.58ms
       { expected_response:true }...: avg=97.21ms min=1.88ms med=58.76ms max=2.21s  p(90)=215.72ms p(95)=272.58ms
   ✓ http_req_failed................: 0.00%   0 out of 11929
     http_req_receiving.............: avg=79.58µs min=17µs   med=72µs    max=1.19ms p(90)=117µs    p(95)=138µs
     http_req_sending...............: avg=29.03µs min=7µs    med=24µs    max=3.95ms p(90)=39µs     p(95)=50µs
     http_req_tls_handshaking.......: avg=0s      min=0s     med=0s      max=0s     p(90)=0s       p(95)=0s
     http_req_waiting...............: avg=97.1ms  min=1.81ms med=58.62ms max=2.21s  p(90)=215.63ms p(95)=272.44ms
     http_reqs......................: 11929   183.504647/s
     iteration_duration.............: avg=97.41ms min=2.02ms med=58.96ms max=2.21s  p(90)=215.89ms p(95)=272.81ms
     iterations.....................: 11929   183.504647/s
     vus............................: 5       min=2              max=20
     vus_max........................: 20      min=20             max=20
```
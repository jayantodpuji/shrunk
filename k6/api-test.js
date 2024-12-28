// import necessary module
import http from 'k6/http';

export default function () {
  // define URL and payload
  const url = 'http://localhost:3002/';
  const payload = JSON.stringify({
   url: 'https://example.com'
  });

  const params = {
    headers: {
      'Content-Type': 'application/json',
    },
  };

  // send a post request and save response as a variable
  const res = http.post(url, payload, params);
  console.log(res)
}

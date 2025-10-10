import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
  stages: [
    { duration: '30s', target: 20 },
    { duration: '1m', target: 50 },
    { duration: '30s', target: 0 },
  ],
};

export default function () {
  const payload = JSON.stringify({
    n: 100,
    min: 1,
    max: 1000,
  });

  const params = {
    headers: {
      'Content-Type': 'application/json',
    },
  };

  const res = http.post('http://localhost:8080/bitonic', payload, params);

  check(res, {
    'status is 200': (r) => r.status === 200,
    'has sequence': (r) => JSON.parse(r.body).sequence !== undefined,
    'is bitonic': (r) => JSON.parse(r.body).is_bitonic === true,
  });

  sleep(1);
}

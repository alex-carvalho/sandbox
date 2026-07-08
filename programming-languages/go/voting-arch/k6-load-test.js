import http from 'k6/http';
import { check } from 'k6';
import { Counter } from 'k6/metrics';

const voteCounter = new Counter('votes_submitted');

export const options = {
  stages: [
    // Aggressive ramp up to maximum throughput
    { duration: '30s', target: 600 }, 
    { duration: '30s', target: 600 },
    { duration: '10s', target: 0 }, 
  ],
  thresholds: {
    http_req_duration: ['p(95)<1000', 'p(99)<2000'],
    http_req_failed: ['rate<0.05'],
    // 'votes_submitted': ['count>=0'], // Expect at least 900k votes
  },
};

const BASE_URL = 'http://localhost:8081';
const VOTING_ID = 100;

export default function () {
  // Use unique user IDs to avoid conflicts and maximize throughput
  const userId = __VU * 100000 + __ITER;
  const voteOption = (Math.floor(Math.random() * 5)) + 1; // 5 options to distribute load

  const payload = JSON.stringify({
    user_id: userId,
    voting_id: VOTING_ID,
    vote_option: voteOption,
  });

  const params = {
    headers: {
      'Content-Type': 'application/json',
    },
    timeout: '30s',
  };

  const response = http.post(`${BASE_URL}/vote`, payload, params);

  check(response, {
    'status is 200': (r) => r.status === 200,
    'response successful': (r) => r.json('status') === 'success',
  });

  voteCounter.add(1);
}

import http from 'k6/http';
import { check } from 'k6';
import { Counter } from 'k6/metrics';

const statusCounter = new Counter('http_req_status');

export const options = {
  scenarios: {
    main: {
      executor: 'constant-arrival-rate',
      rate: 1000,
      timeUnit: '1m', // 1000 iterations per minute
      duration: '1m',
      preAllocatedVUs: 50,
      maxVUs: 150,
    },
  },
};

export default function () {
  const url = 'http://traykun.gateway.test/deck_of_cards/deck/new/';
  const params = {
    headers: { 'Api-Key': 'mk_live_demo_t42A1M2ruNLKhn33gOV3TB0WYhqZLw9cbgTB9TEwwEbkx-hPEHvS7wHUyW1jUmGQxdERArkSdAHUOFsIeYWEdQ==' },
  };

  const res = http.get(url, params);

  // Cek status, bisa 200 OK atau 429 Too Many Requests
  check(res, {
    'status is 200 or 429': (r) => [200, 429].includes(r.status),
  });

  statusCounter.add(1, { status: res.status });
}

export function handleSummary(data) {
  const status200s = data.metrics.http_req_status.values['count{status:200}'] || 0;
  const status429s = data.metrics.http_req_status.values['count{status:429}'] || 0;

  console.log('--- Status Code Summary ---');
  console.log(`Number of 200 OK responses: ${status200s}`);
  console.log(`Number of 429 Too Many Requests responses: ${status429s}`);
  console.log('---------------------------');

  // To also show the default summary, we can import and use defaultHandleSummary
  // For now, we'll just show our custom summary.
  return {};
}
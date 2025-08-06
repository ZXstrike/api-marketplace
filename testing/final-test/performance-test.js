import http from 'k6/http';
import { check } from 'k6';

export const options = {
  scenarios: {
    main: {
      executor: 'constant-arrival-rate',
      rate: 1500,
      timeUnit: '1m', // 1000 iterations per minute
      duration: '1m',
      preAllocatedVUs: 50,
      maxVUs: 200,
    },
  },
};

export default function () {
  const url = 'http://traykun.api.gateway.test/deck_of_cards/deck/new/';
  const params = {
    headers: { 'Api-Key': 'mk_live_demo_t42A1M2ruNLKhn33gOV3TB0WYhqZLw9cbgTB9TEwwEbkx-hPEHvS7wHUyW1jUmGQxdERArkSdAHUOFsIeYWEdQ==' },
  };

  const res = http.get(url, params);

  // Cek status, bisa 200 OK atau 429 Too Many Requests
  check(res, {
    'status is 200 or 429': (r) => [200, 429].includes(r.status),
  });
}
import http from 'k6/http';
import { check } from 'k6';

export const options = {
  // 10 virtual user, masing-masing mencoba 1 request secara bersamaan
  // Total ada 10 permintaan, padahal saldo hanya cukup untuk 5
  vus: 20,
  iterations: 60,
};

export default function () {
  const url = 'http://traykun.api.gateway.test/deck_of_cards/deck/new/';
  const params = {
    headers: { 'Api-Key': 'mk_live_demo_t42A1M2ruNLKhn33gOV3TB0WYhqZLw9cbgTB9TEwwEbkx-hPEHvS7wHUyW1jUmGQxdERArkSdAHUOFsIeYWEdQ==' },
  };

  const res = http.get(url, params);

  // Cek status, bisa 200 OK atau 402 Payment Required
  check(res, {
    'status is 200 or 402': (r) => [200, 402].includes(r.status),
  });
}
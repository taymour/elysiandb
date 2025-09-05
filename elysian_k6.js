import http from 'k6/http';
import { check } from 'k6';

const BASE = __ENV.BASE_URL || 'http://localhost:8089';
const KEYS = parseInt(__ENV.KEYS || '5000', 10);
const VUS  = parseInt(__ENV.VUS  || '200', 10);
const DURATION = __ENV.DURATION || '30s';

const pad = (n) => n.toString().padStart(6, '0');

function keyForThisVU() {
  const perVu = Math.max(1, Math.floor(KEYS / VUS));
  const base = ((__VU - 1) * perVu);
  const idx = base + (__ITER % perVu);
  const id = Math.min(idx + 1, KEYS);
  return `bench${pad(id)}`;
}

export const options = {
  vus: VUS,
  duration: DURATION,
  thresholds: {
    http_req_failed:   ['rate<0.01'],
    http_req_duration: ['p(95)<25'],
  },
};

export default function () {
  const key = keyForThisVU();
  const useTTL = (__ITER % 2) === 0; // 1 PUT sur 2 avec TTL
  const urlPut = useTTL ? `${BASE}/kv/${key}?ttl=50` : `${BASE}/kv/${key}`;

  const put = http.put(
    urlPut,
    `value=${__ITER}`,
    { headers: { 'Content-Type': 'application/x-www-form-urlencoded' }, tags: { name: 'kv_put' } }
  );
  check(put, { 'PUT 204': (r) => r.status === 204 });

  // GET
  const get = http.get(`${BASE}/kv/${key}`, { tags: { name: 'kv_get' } });
  check(get, { 'GET 200': (r) => r.status === 200 });

  // DELETE
  const del = http.del(`${BASE}/kv/${key}`, null, { tags: { name: 'kv_del' } });
  check(del, { 'DEL 204': (r) => r.status === 204 });
}

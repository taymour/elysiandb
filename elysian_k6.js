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
  const val = `value=${__ITER}`;
  const useTTL = (__ITER % 2) === 0;
  const urlPut = useTTL ? `${BASE}/kv/${key}?ttl=50` : `${BASE}/kv/${key}`;

  const put = http.put(
    urlPut,
    val,
    { headers: { 'Content-Type': 'application/x-www-form-urlencoded' }, tags: { name: 'kv_put' } }
  );
  check(put, { 'PUT 204': (r) => r.status === 204 });

  const get = http.get(`${BASE}/kv/${key}`, { tags: { name: 'kv_get' } });
  const okGet = check(get, {
    'GET 200': (r) => r.status === 200,
    'GET JSON has key/value': (r) => {
      try {
        const body = JSON.parse(r.body);
        return body && body.key === key && body.value === val;
      } catch (_) {
        return false;
      }
    },
  });

  if ((__ITER % 10) === 0) {
    const missKey = `${key}-missing`;
    const mget = http.get(`${BASE}/kv/mget?keys=${encodeURIComponent(key)},${encodeURIComponent(missKey)}`, { tags: { name: 'kv_mget' } });
    check(mget, {
      'MGET 200': (r) => r.status === 200,
      'MGET JSON array ok': (r) => {
        try {
          const arr = JSON.parse(r.body);
          if (!Array.isArray(arr) || arr.length !== 2) return false;
          const a = arr[0], b = arr[1];
          return a.key === key && a.value === val && b.key === missKey && b.value === null;
        } catch (_) {
          return false;
        }
      },
    });
  }

  const del = http.del(`${BASE}/kv/${key}`, null, { tags: { name: 'kv_del' } });
  check(del, { 'DEL 204': (r) => r.status === 204 });

  if ((__ITER % 5) === 0) {
    const getAfterDel = http.get(`${BASE}/kv/${key}`, { tags: { name: 'kv_get_after_del' } });
    check(getAfterDel, {
      'GET after DEL 404': (r) => r.status === 404,
      'GET after DEL JSON null': (r) => {
        try {
          const body = JSON.parse(r.body);
          return body && body.key === key && body.value === null;
        } catch (_) {
          return false;
        }
      },
    });
  }
}

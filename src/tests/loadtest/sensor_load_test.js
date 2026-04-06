import { group, sleep } from 'k6';
import {
  buildValidEventPayload,
  checkEventAccepted,
  postEvent,
} from './sensor_shared.js';

export const options = {
  stages: [
    { duration: '20s', target: 12 },
    { duration: '45s', target: 12 },
    { duration: '20s', target: 28 },
    { duration: '45s', target: 28 },
    { duration: '20s', target: 0 },
  ],
  thresholds: {
    http_req_failed: ['rate<0.01'],
    http_req_duration: ['p(95)<800'],
    sensor_success_rate: ['rate>0.97'],
  },
};

export default function () {
  group('POST /events em carga normal', () => {
    const res = postEvent(buildValidEventPayload());
    checkEventAccepted(res);
  });

  sleep(0.08);
}

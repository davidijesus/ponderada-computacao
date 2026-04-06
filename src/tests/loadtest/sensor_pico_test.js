import { group, sleep } from 'k6';
import {
  buildValidEventPayload,
  checkEventAccepted,
  postEvent,
} from './sensor_shared.js';

export const options = {
  stages: [
    { duration: '15s', target: 40 },
    { duration: '10s', target: 160 },
    { duration: '30s', target: 160 },
    { duration: '10s', target: 30 },
    { duration: '15s', target: 0 },
  ],
  thresholds: {
    http_req_failed: ['rate<0.01'],
    http_req_duration: ['p(95)<1500'],
    sensor_success_rate: ['rate>0.95'],
  },
};

export default function () {
  group('POST /events em pico', () => {
    const res = postEvent(buildValidEventPayload());
    checkEventAccepted(res);
  });

  sleep(0.04);
}

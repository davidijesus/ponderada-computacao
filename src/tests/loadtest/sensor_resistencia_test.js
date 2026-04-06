import { group, sleep } from 'k6';
import {
  buildValidEventPayload,
  checkEventAccepted,
  postEvent,
} from './sensor_shared.js';

export const options = {
  stages: [
    { duration: '30s', target: 25 },
    { duration: '40s', target: 45 },
    { duration: '3m30s', target: 45 },
    { duration: '30s', target: 0 },
  ],
  thresholds: {
    http_req_failed: ['rate<0.01'],
    http_req_duration: ['p(95)<1000'],
    sensor_success_rate: ['rate>0.97'],
  },
};

export default function () {
  group('POST /events em resistência', () => {
    const res = postEvent(buildValidEventPayload());
    checkEventAccepted(res);
  });

  sleep(0.08);
}

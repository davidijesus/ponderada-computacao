import { group, sleep } from 'k6';
import {
  buildValidEventPayload,
  checkEventAccepted,
  postEvent,
} from './sensor_shared.js';

export const options = {
  stages: [
    { duration: '20s', target: 18 },
    { duration: '40s', target: 35 },
    { duration: '50s', target: 55 },
    { duration: '20s', target: 0 },
  ],
  thresholds: {
    http_req_failed: ['rate<0.01'],
    http_req_duration: ['p(95)<1200'],
    sensor_success_rate: ['rate>0.96'],
  },
};

export default function () {
  group('POST /events sob stress', () => {
    const res = postEvent(buildValidEventPayload());
    checkEventAccepted(res);
  });

  sleep(0.04);
}

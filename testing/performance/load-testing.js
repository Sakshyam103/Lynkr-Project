/**
 * Performance Load Testing
 * Load tests for high-traffic scenarios during events
 */

import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate } from 'k6/metrics';

// Custom metrics
const errorRate = new Rate('errors');

// Test configuration
export const options = {
  stages: [
    { duration: '2m', target: 100 }, // Ramp up to 100 users
    { duration: '5m', target: 100 }, // Stay at 100 users
    { duration: '2m', target: 200 }, // Ramp up to 200 users
    { duration: '5m', target: 200 }, // Stay at 200 users
    { duration: '2m', target: 500 }, // Ramp up to 500 users
    { duration: '10m', target: 500 }, // Stay at 500 users (peak event load)
    { duration: '5m', target: 0 }, // Ramp down
  ],
  thresholds: {
    http_req_duration: ['p(95)<500'], // 95% of requests under 500ms
    http_req_failed: ['rate<0.1'], // Error rate under 10%
    errors: ['rate<0.1'],
  },
};

// Test data
const BASE_URL = 'https://api.lynkr.com/v1';
let authTokens = [];
let eventIds = [];

export function setup() {
  // Create test users and events for load testing
  const users = [];
  const events = [];
  
  // Create 50 test users
  for (let i = 0; i < 50; i++) {
    const response = http.post(`${BASE_URL}/users/register`, JSON.stringify({
      email: `loadtest${i}@example.com`,
      password: 'LoadTest123!',
      name: `Load Test User ${i}`
    }), {
      headers: { 'Content-Type': 'application/json' }
    });
    
    if (response.status === 201) {
      users.push(response.json().token);
    }
  }
  
  // Create test events
  const brandToken = 'brand_test_token'; // Would be obtained from brand login
  for (let i = 0; i < 10; i++) {
    const response = http.post(`${BASE_URL}/events`, JSON.stringify({
      name: `Load Test Event ${i}`,
      description: 'Load testing event',
      startDate: '2024-12-31T10:00:00Z',
      endDate: '2024-12-31T18:00:00Z',
      location: {
        latitude: 37.7749 + (Math.random() * 0.01),
        longitude: -122.4194 + (Math.random() * 0.01),
        address: `${100 + i} Test St, San Francisco, CA`
      }
    }), {
      headers: { 
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${brandToken}`
      }
    });
    
    if (response.status === 201) {
      events.push(response.json().id);
    }
  }
  
  return { users, events };
}

export default function(data) {
  const userToken = data.users[Math.floor(Math.random() * data.users.length)];
  const eventId = data.events[Math.floor(Math.random() * data.events.length)];
  
  // Simulate user behavior during event
  eventDiscovery(userToken);
  sleep(1);
  
  eventCheckIn(userToken, eventId);
  sleep(2);
  
  contentCreation(userToken, eventId);
  sleep(1);
  
  pollParticipation(userToken, eventId);
  sleep(1);
  
  analyticsView(userToken, eventId);
}

function eventDiscovery(token) {
  const response = http.get(`${BASE_URL}/events`, {
    headers: { 'Authorization': `Bearer ${token}` }
  });
  
  const success = check(response, {
    'event discovery status is 200': (r) => r.status === 200,
    'event discovery response time < 200ms': (r) => r.timings.duration < 200,
    'events returned': (r) => JSON.parse(r.body).length > 0,
  });
  
  errorRate.add(!success);
}

function eventCheckIn(token, eventId) {
  const response = http.post(`${BASE_URL}/events/${eventId}/checkin`, JSON.stringify({
    latitude: 37.7749 + (Math.random() * 0.001),
    longitude: -122.4194 + (Math.random() * 0.001)
  }), {
    headers: { 
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token}`
    }
  });
  
  const success = check(response, {
    'check-in status is 200': (r) => r.status === 200,
    'check-in response time < 300ms': (r) => r.timings.duration < 300,
    'check-in has timestamp': (r) => JSON.parse(r.body).checkinTime !== undefined,
  });
  
  errorRate.add(!success);
}

function contentCreation(token, eventId) {
  const response = http.post(`${BASE_URL}/content`, JSON.stringify({
    eventId: eventId,
    mediaType: 'photo',
    caption: `Load test content #event${eventId}`,
    tags: ['loadtest', 'performance']
  }), {
    headers: { 
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token}`
    }
  });
  
  const success = check(response, {
    'content creation status is 201': (r) => r.status === 201,
    'content creation response time < 500ms': (r) => r.timings.duration < 500,
    'content has ID': (r) => JSON.parse(r.body).id !== undefined,
  });
  
  errorRate.add(!success);
}

function pollParticipation(token, eventId) {
  // First get available polls
  const pollsResponse = http.get(`${BASE_URL}/events/${eventId}/polls`, {
    headers: { 'Authorization': `Bearer ${token}` }
  });
  
  if (pollsResponse.status === 200) {
    const polls = JSON.parse(pollsResponse.body);
    if (polls.length > 0) {
      const poll = polls[0];
      const option = poll.options[Math.floor(Math.random() * poll.options.length)];
      
      const voteResponse = http.post(`${BASE_URL}/feedback/polls/vote`, JSON.stringify({
        pollId: poll.id,
        optionId: option.id
      }), {
        headers: { 
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`
        }
      });
      
      const success = check(voteResponse, {
        'poll vote status is 200': (r) => r.status === 200,
        'poll vote response time < 200ms': (r) => r.timings.duration < 200,
      });
      
      errorRate.add(!success);
    }
  }
}

function analyticsView(token, eventId) {
  const response = http.get(`${BASE_URL}/events/${eventId}/analytics/realtime`, {
    headers: { 'Authorization': `Bearer ${token}` }
  });
  
  const success = check(response, {
    'analytics status is 200': (r) => r.status === 200,
    'analytics response time < 300ms': (r) => r.timings.duration < 300,
    'analytics has metrics': (r) => Object.keys(JSON.parse(r.body)).length > 0,
  });
  
  errorRate.add(!success);
}

// Brand portal load testing
export function brandPortalLoad() {
  const brandToken = 'brand_test_token';
  
  // Dashboard load
  const dashboardResponse = http.get(`${BASE_URL}/brands/dashboard`, {
    headers: { 'Authorization': `Bearer ${brandToken}` }
  });
  
  check(dashboardResponse, {
    'brand dashboard status is 200': (r) => r.status === 200,
    'brand dashboard response time < 400ms': (r) => r.timings.duration < 400,
  });
  
  // Real-time analytics
  const analyticsResponse = http.get(`${BASE_URL}/events/test-event/analytics/realtime`, {
    headers: { 'Authorization': `Bearer ${brandToken}` }
  });
  
  check(analyticsResponse, {
    'brand analytics status is 200': (r) => r.status === 200,
    'brand analytics response time < 500ms': (r) => r.timings.duration < 500,
  });
}

// Database performance testing
export function databaseLoad() {
  const token = 'test_token';
  
  // Heavy query load
  const queries = [
    `${BASE_URL}/events?limit=100`,
    `${BASE_URL}/content?limit=50`,
    `${BASE_URL}/analytics/events/summary`,
    `${BASE_URL}/users/rewards`,
    `${BASE_URL}/events/test-event/conversion-funnel`
  ];
  
  queries.forEach(url => {
    const response = http.get(url, {
      headers: { 'Authorization': `Bearer ${token}` }
    });
    
    check(response, {
      [`${url} status is 200`]: (r) => r.status === 200,
      [`${url} response time < 1000ms`]: (r) => r.timings.duration < 1000,
    });
  });
}

// Memory and resource testing
export function resourceUsageTest() {
  const token = 'test_token';
  
  // Create large payload
  const largeContent = {
    eventId: 'test-event',
    mediaType: 'video',
    caption: 'A'.repeat(1000), // Large caption
    tags: Array(100).fill('tag').map((t, i) => `${t}${i}`), // Many tags
    metadata: {
      duration: 120,
      size: 50000000, // 50MB
      resolution: '1920x1080'
    }
  };
  
  const response = http.post(`${BASE_URL}/content`, JSON.stringify(largeContent), {
    headers: { 
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token}`
    }
  });
  
  check(response, {
    'large content creation handles properly': (r) => r.status === 201 || r.status === 413,
    'large content response time reasonable': (r) => r.timings.duration < 5000,
  });
}

export function teardown(data) {
  // Cleanup test data
  console.log('Load test completed. Cleaning up test data...');
  
  // In a real implementation, would clean up test users and events
  // This helps prevent database bloat from repeated load testing
}
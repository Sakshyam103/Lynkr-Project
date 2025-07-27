/**
 * API Integration Tests
 * Tests for API endpoints and service integrations
 */

const request = require('supertest');
const app = require('../../backend/cmd/api/main');

describe('API Integration Tests', () => {
  let authToken;
  let brandToken;
  let testEventId;
  let testUserId;

  beforeAll(async () => {
    // Setup test data
    const userResponse = await request(app)
      .post('/api/v1/users/register')
      .send({
        email: 'test@example.com',
        password: 'SecurePass123!',
        name: 'Test User'
      });
    
    authToken = userResponse.body.token;
    testUserId = userResponse.body.user.id;

    const brandResponse = await request(app)
      .post('/api/v1/brands/login')
      .send({
        email: 'brand@example.com',
        password: 'BrandPass123!'
      });
    
    brandToken = brandResponse.body.token;
  });

  describe('User Authentication', () => {
    test('POST /api/v1/users/register - should create new user', async () => {
      const response = await request(app)
        .post('/api/v1/users/register')
        .send({
          email: 'newuser@example.com',
          password: 'SecurePass123!',
          name: 'New User'
        });

      expect(response.status).toBe(201);
      expect(response.body).toHaveProperty('token');
      expect(response.body.user).toHaveProperty('id');
      expect(response.body.user.email).toBe('newuser@example.com');
    });

    test('POST /api/v1/users/login - should authenticate user', async () => {
      const response = await request(app)
        .post('/api/v1/users/login')
        .send({
          email: 'test@example.com',
          password: 'SecurePass123!'
        });

      expect(response.status).toBe(200);
      expect(response.body).toHaveProperty('token');
      expect(response.body.user.email).toBe('test@example.com');
    });

    test('POST /api/v1/users/login - should reject invalid credentials', async () => {
      const response = await request(app)
        .post('/api/v1/users/login')
        .send({
          email: 'test@example.com',
          password: 'wrongpassword'
        });

      expect(response.status).toBe(401);
      expect(response.body).toHaveProperty('error');
    });
  });

  describe('Event Management', () => {
    test('POST /api/v1/events - should create event', async () => {
      const response = await request(app)
        .post('/api/v1/events')
        .set('Authorization', `Bearer ${brandToken}`)
        .send({
          name: 'Test Event',
          description: 'Integration test event',
          startDate: '2024-12-31T10:00:00Z',
          endDate: '2024-12-31T18:00:00Z',
          location: {
            latitude: 37.7749,
            longitude: -122.4194,
            address: '123 Test St, San Francisco, CA'
          }
        });

      expect(response.status).toBe(201);
      expect(response.body).toHaveProperty('id');
      expect(response.body.name).toBe('Test Event');
      testEventId = response.body.id;
    });

    test('GET /api/v1/events/:id - should get event details', async () => {
      const response = await request(app)
        .get(`/api/v1/events/${testEventId}`)
        .set('Authorization', `Bearer ${authToken}`);

      expect(response.status).toBe(200);
      expect(response.body.id).toBe(testEventId);
      expect(response.body.name).toBe('Test Event');
    });

    test('POST /api/v1/events/:id/checkin - should check in user', async () => {
      const response = await request(app)
        .post(`/api/v1/events/${testEventId}/checkin`)
        .set('Authorization', `Bearer ${authToken}`)
        .send({
          latitude: 37.7749,
          longitude: -122.4194
        });

      expect(response.status).toBe(200);
      expect(response.body).toHaveProperty('checkinTime');
    });
  });

  describe('Content Management', () => {
    let testContentId;

    test('POST /api/v1/content - should create content', async () => {
      const response = await request(app)
        .post('/api/v1/content')
        .set('Authorization', `Bearer ${authToken}`)
        .send({
          eventId: testEventId,
          mediaType: 'photo',
          caption: 'Test content #integration',
          tags: ['integration', 'test']
        });

      expect(response.status).toBe(201);
      expect(response.body).toHaveProperty('id');
      expect(response.body.caption).toBe('Test content #integration');
      testContentId = response.body.id;
    });

    test('GET /api/v1/content/:id - should get content details', async () => {
      const response = await request(app)
        .get(`/api/v1/content/${testContentId}`)
        .set('Authorization', `Bearer ${authToken}`);

      expect(response.status).toBe(200);
      expect(response.body.id).toBe(testContentId);
      expect(response.body.caption).toBe('Test content #integration');
    });

    test('PUT /api/v1/content/:id/permissions - should update permissions', async () => {
      const response = await request(app)
        .put(`/api/v1/content/${testContentId}/permissions`)
        .set('Authorization', `Bearer ${authToken}`)
        .send({
          brandUsage: true,
          publicDisplay: false,
          attribution: 'required'
        });

      expect(response.status).toBe(200);
      expect(response.body.permissions.brandUsage).toBe(true);
    });
  });

  describe('Analytics and Reporting', () => {
    test('GET /api/v1/events/:id/analytics/attendance - should get attendance analytics', async () => {
      const response = await request(app)
        .get(`/api/v1/events/${testEventId}/analytics/attendance`)
        .set('Authorization', `Bearer ${brandToken}`);

      expect(response.status).toBe(200);
      expect(response.body).toHaveProperty('totalAttendees');
      expect(response.body).toHaveProperty('uniqueUsers');
    });

    test('GET /api/v1/events/:id/analytics/engagement - should get engagement metrics', async () => {
      const response = await request(app)
        .get(`/api/v1/events/${testEventId}/analytics/engagement`)
        .set('Authorization', `Bearer ${brandToken}`);

      expect(response.status).toBe(200);
      expect(response.body).toHaveProperty('engagementRate');
      expect(response.body).toHaveProperty('contentCreated');
    });

    test('GET /api/v1/events/:id/conversion-funnel - should get conversion funnel', async () => {
      const response = await request(app)
        .get(`/api/v1/events/${testEventId}/conversion-funnel`)
        .set('Authorization', `Bearer ${brandToken}`);

      expect(response.status).toBe(200);
      expect(response.body).toHaveProperty('stages');
      expect(response.body).toHaveProperty('totalUsers');
      expect(response.body).toHaveProperty('roi');
    });
  });

  describe('E-commerce Integration', () => {
    test('POST /api/v1/ecommerce/integrations - should create integration', async () => {
      const response = await request(app)
        .post('/api/v1/ecommerce/integrations')
        .set('Authorization', `Bearer ${brandToken}`)
        .send({
          platform: 'shopify',
          storeUrl: 'test-store.myshopify.com',
          accessToken: 'test_token'
        });

      expect(response.status).toBe(201);
      expect(response.body).toHaveProperty('id');
      expect(response.body.platform).toBe('shopify');
    });

    test('POST /api/v1/ecommerce/purchases - should track purchase', async () => {
      const response = await request(app)
        .post('/api/v1/ecommerce/purchases')
        .set('Authorization', `Bearer ${authToken}`)
        .send({
          orderId: 'order_123',
          eventId: testEventId,
          amount: 99.99,
          currency: 'USD',
          products: [{
            id: 'product_123',
            name: 'Test Product',
            price: 99.99,
            quantity: 1
          }]
        });

      expect(response.status).toBe(201);
      expect(response.body).toHaveProperty('id');
      expect(response.body.amount).toBe(99.99);
    });
  });

  describe('Discount Codes', () => {
    let testCodeId;

    test('POST /api/v1/discount/generate - should generate discount code', async () => {
      const response = await request(app)
        .post('/api/v1/discount/generate')
        .set('Authorization', `Bearer ${brandToken}`)
        .send({
          eventId: testEventId,
          discountType: 'percentage',
          discountValue: 20,
          maxUses: 100
        });

      expect(response.status).toBe(201);
      expect(response.body).toHaveProperty('code');
      expect(response.body.discountValue).toBe(20);
      testCodeId = response.body.id;
    });

    test('GET /api/v1/discount/codes/:code/validate - should validate code', async () => {
      const codeResponse = await request(app)
        .get(`/api/v1/discount/codes/${testCodeId}`)
        .set('Authorization', `Bearer ${brandToken}`);
      
      const code = codeResponse.body.code;

      const response = await request(app)
        .get(`/api/v1/discount/codes/${code}/validate`)
        .set('Authorization', `Bearer ${authToken}`);

      expect(response.status).toBe(200);
      expect(response.body.valid).toBe(true);
      expect(response.body.discountValue).toBe(20);
    });
  });

  describe('Security and Privacy', () => {
    test('POST /api/v1/security/scan - should run security scan', async () => {
      const response = await request(app)
        .post('/api/v1/security/scan')
        .set('Authorization', `Bearer ${brandToken}`);

      expect(response.status).toBe(200);
      expect(response.body).toHaveProperty('vulnerabilities');
      expect(response.body.scan_completed).toBe(true);
    });

    test('GET /api/v1/users/data/export - should export user data', async () => {
      const response = await request(app)
        .get('/api/v1/users/data/export')
        .set('Authorization', `Bearer ${authToken}`);

      expect(response.status).toBe(200);
      expect(response.body).toHaveProperty('profile');
      expect(response.body).toHaveProperty('attendances');
      expect(response.body).toHaveProperty('content');
    });

    test('POST /api/v1/users/data/delete - should schedule data deletion', async () => {
      const response = await request(app)
        .post('/api/v1/users/data/delete')
        .set('Authorization', `Bearer ${authToken}`);

      expect(response.status).toBe(200);
      expect(response.body.status).toBe('deletion_scheduled');
    });
  });

  describe('Rate Limiting', () => {
    test('should enforce rate limits', async () => {
      const requests = [];
      
      // Make 101 requests rapidly (assuming 100 req/min limit)
      for (let i = 0; i < 101; i++) {
        requests.push(
          request(app)
            .get('/api/v1/events')
            .set('Authorization', `Bearer ${authToken}`)
        );
      }

      const responses = await Promise.all(requests);
      const rateLimitedResponses = responses.filter(r => r.status === 429);
      
      expect(rateLimitedResponses.length).toBeGreaterThan(0);
    });
  });

  describe('Error Handling', () => {
    test('should handle invalid JSON', async () => {
      const response = await request(app)
        .post('/api/v1/events')
        .set('Authorization', `Bearer ${brandToken}`)
        .set('Content-Type', 'application/json')
        .send('invalid json');

      expect(response.status).toBe(400);
      expect(response.body).toHaveProperty('error');
    });

    test('should handle missing authentication', async () => {
      const response = await request(app)
        .get('/api/v1/events');

      expect(response.status).toBe(401);
      expect(response.body).toHaveProperty('error');
    });

    test('should handle resource not found', async () => {
      const response = await request(app)
        .get('/api/v1/events/nonexistent')
        .set('Authorization', `Bearer ${authToken}`);

      expect(response.status).toBe(404);
      expect(response.body).toHaveProperty('error');
    });
  });
});
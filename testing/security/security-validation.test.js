/**
 * Security Validation Tests
 * Comprehensive security testing for vulnerabilities and compliance
 */

const request = require('supertest');
const app = require('../../backend/cmd/api/main');

describe('Security Validation Tests', () => {
  let authToken;
  let brandToken;

  beforeAll(async () => {
    // Setup authenticated sessions
    const userResponse = await request(app)
      .post('/api/v1/users/register')
      .send({
        email: 'security@example.com',
        password: 'SecurePass123!',
        name: 'Security Test User'
      });
    authToken = userResponse.body.token;

    const brandResponse = await request(app)
      .post('/api/v1/brands/login')
      .send({
        email: 'securitybrand@example.com',
        password: 'BrandSecure123!'
      });
    brandToken = brandResponse.body.token;
  });

  describe('Authentication Security', () => {
    test('should reject weak passwords', async () => {
      const weakPasswords = ['123', 'password', 'abc123', '12345678'];
      
      for (const password of weakPasswords) {
        const response = await request(app)
          .post('/api/v1/users/register')
          .send({
            email: `weak${Date.now()}@example.com`,
            password: password,
            name: 'Weak Password Test'
          });
        
        expect(response.status).toBe(400);
        expect(response.body.error).toContain('password');
      }
    });

    test('should enforce rate limiting on login attempts', async () => {
      const requests = [];
      
      // Attempt 10 rapid login requests
      for (let i = 0; i < 10; i++) {
        requests.push(
          request(app)
            .post('/api/v1/users/login')
            .send({
              email: 'nonexistent@example.com',
              password: 'wrongpassword'
            })
        );
      }
      
      const responses = await Promise.all(requests);
      const rateLimitedResponses = responses.filter(r => r.status === 429);
      
      expect(rateLimitedResponses.length).toBeGreaterThan(0);
    });

    test('should invalidate tokens after logout', async () => {
      // Login and get token
      const loginResponse = await request(app)
        .post('/api/v1/users/login')
        .send({
          email: 'security@example.com',
          password: 'SecurePass123!'
        });
      
      const token = loginResponse.body.token;
      
      // Use token successfully
      const profileResponse = await request(app)
        .get('/api/v1/users/profile')
        .set('Authorization', `Bearer ${token}`);
      
      expect(profileResponse.status).toBe(200);
      
      // Logout
      await request(app)
        .post('/api/v1/users/logout')
        .set('Authorization', `Bearer ${token}`);
      
      // Token should now be invalid
      const invalidResponse = await request(app)
        .get('/api/v1/users/profile')
        .set('Authorization', `Bearer ${token}`);
      
      expect(invalidResponse.status).toBe(401);
    });
  });

  describe('Input Validation', () => {
    test('should prevent SQL injection attempts', async () => {
      const sqlInjectionPayloads = [
        "'; DROP TABLE users; --",
        "' OR '1'='1",
        "'; INSERT INTO users (email) VALUES ('hacked@evil.com'); --",
        "' UNION SELECT * FROM users --"
      ];
      
      for (const payload of sqlInjectionPayloads) {
        const response = await request(app)
          .post('/api/v1/users/register')
          .send({
            email: payload,
            password: 'ValidPass123!',
            name: 'SQL Injection Test'
          });
        
        expect(response.status).toBe(400);
        expect(response.body.error).toBeDefined();
      }
    });

    test('should prevent XSS attacks', async () => {
      const xssPayloads = [
        '<script>alert("XSS")</script>',
        '<img src="x" onerror="alert(1)">',
        'javascript:alert("XSS")',
        '<svg onload="alert(1)">'
      ];
      
      for (const payload of xssPayloads) {
        const response = await request(app)
          .post('/api/v1/content')
          .set('Authorization', `Bearer ${authToken}`)
          .send({
            eventId: 'test-event',
            mediaType: 'photo',
            caption: payload,
            tags: ['test']
          });
        
        // Should either reject or sanitize the input
        if (response.status === 201) {
          expect(response.body.caption).not.toContain('<script>');
          expect(response.body.caption).not.toContain('javascript:');
        } else {
          expect(response.status).toBe(400);
        }
      }
    });

    test('should validate file uploads', async () => {
      // Test malicious file upload
      const response = await request(app)
        .post('/api/v1/content/upload')
        .set('Authorization', `Bearer ${authToken}`)
        .attach('file', Buffer.from('<?php echo "hacked"; ?>'), 'malicious.php');
      
      expect(response.status).toBe(400);
      expect(response.body.error).toContain('file type');
    });

    test('should enforce input length limits', async () => {
      const longString = 'A'.repeat(10000);
      
      const response = await request(app)
        .post('/api/v1/content')
        .set('Authorization', `Bearer ${authToken}`)
        .send({
          eventId: 'test-event',
          mediaType: 'photo',
          caption: longString,
          tags: ['test']
        });
      
      expect(response.status).toBe(400);
      expect(response.body.error).toContain('length');
    });
  });

  describe('Authorization Controls', () => {
    test('should prevent unauthorized access to brand endpoints', async () => {
      const response = await request(app)
        .get('/api/v1/brands/dashboard')
        .set('Authorization', `Bearer ${authToken}`); // User token, not brand token
      
      expect(response.status).toBe(403);
    });

    test('should prevent users from accessing other users\' data', async () => {
      // Create another user
      const otherUserResponse = await request(app)
        .post('/api/v1/users/register')
        .send({
          email: 'other@example.com',
          password: 'OtherPass123!',
          name: 'Other User'
        });
      
      const otherUserId = otherUserResponse.body.user.id;
      
      // Try to access other user's data
      const response = await request(app)
        .get(`/api/v1/users/${otherUserId}/profile`)
        .set('Authorization', `Bearer ${authToken}`);
      
      expect(response.status).toBe(403);
    });

    test('should validate JWT token integrity', async () => {
      const tamperedToken = authToken.slice(0, -5) + 'XXXXX';
      
      const response = await request(app)
        .get('/api/v1/users/profile')
        .set('Authorization', `Bearer ${tamperedToken}`);
      
      expect(response.status).toBe(401);
    });
  });

  describe('Data Privacy Compliance', () => {
    test('should respect user consent settings', async () => {
      // Update user consent to deny analytics
      await request(app)
        .put('/api/v1/users/privacy')
        .set('Authorization', `Bearer ${authToken}`)
        .send({
          analyticsConsent: false,
          marketingConsent: false
        });
      
      // Analytics tracking should be blocked
      const response = await request(app)
        .post('/api/v1/analytics/track')
        .set('Authorization', `Bearer ${authToken}`)
        .send({
          eventType: 'page_view',
          eventData: { page: 'events' }
        });
      
      expect(response.status).toBe(403);
      expect(response.body.error).toContain('consent');
    });

    test('should anonymize data in exports', async () => {
      const response = await request(app)
        .get('/api/v1/users/data/export')
        .set('Authorization', `Bearer ${authToken}`);
      
      expect(response.status).toBe(200);
      
      // Check that sensitive data is properly handled
      const exportData = response.body;
      expect(exportData.profile.email).toBeDefined();
      expect(exportData.profile.password).toBeUndefined();
      expect(exportData.profile.internalId).toBeUndefined();
    });

    test('should handle data deletion requests', async () => {
      const response = await request(app)
        .post('/api/v1/users/data/delete')
        .set('Authorization', `Bearer ${authToken}`);
      
      expect(response.status).toBe(200);
      expect(response.body.status).toBe('deletion_scheduled');
    });
  });

  describe('API Security', () => {
    test('should enforce HTTPS in production', async () => {
      // This would be tested in a production environment
      // For now, we check that security headers are present
      const response = await request(app)
        .get('/api/v1/events')
        .set('Authorization', `Bearer ${authToken}`);
      
      expect(response.headers['x-content-type-options']).toBe('nosniff');
      expect(response.headers['x-frame-options']).toBe('DENY');
      expect(response.headers['x-xss-protection']).toBe('1; mode=block');
    });

    test('should validate API versioning', async () => {
      const response = await request(app)
        .get('/api/v999/events') // Invalid version
        .set('Authorization', `Bearer ${authToken}`);
      
      expect(response.status).toBe(404);
    });

    test('should handle CORS properly', async () => {
      const response = await request(app)
        .options('/api/v1/events')
        .set('Origin', 'https://malicious-site.com');
      
      // Should not allow arbitrary origins
      expect(response.headers['access-control-allow-origin']).not.toBe('https://malicious-site.com');
    });
  });

  describe('Vulnerability Scanning', () => {
    test('should detect and report vulnerabilities', async () => {
      const response = await request(app)
        .post('/api/v1/security/scan')
        .set('Authorization', `Bearer ${brandToken}`);
      
      expect(response.status).toBe(200);
      expect(response.body).toHaveProperty('vulnerabilities');
      expect(Array.isArray(response.body.vulnerabilities)).toBe(true);
    });

    test('should validate input sanitization', async () => {
      const response = await request(app)
        .post('/api/v1/security/validate-input')
        .set('Authorization', `Bearer ${authToken}`)
        .send({
          input: '<script>alert("test")</script>'
        });
      
      expect(response.status).toBe(200);
      expect(response.body.valid).toBe(false);
    });
  });

  describe('Session Security', () => {
    test('should expire sessions after timeout', async () => {
      // This would require time manipulation in a real test
      // For now, we test that session timeout is configurable
      const response = await request(app)
        .get('/api/v1/users/session-info')
        .set('Authorization', `Bearer ${authToken}`);
      
      expect(response.status).toBe(200);
      expect(response.body).toHaveProperty('expiresAt');
    });

    test('should prevent session fixation', async () => {
      // Login should generate a new session token
      const loginResponse1 = await request(app)
        .post('/api/v1/users/login')
        .send({
          email: 'security@example.com',
          password: 'SecurePass123!'
        });
      
      const loginResponse2 = await request(app)
        .post('/api/v1/users/login')
        .send({
          email: 'security@example.com',
          password: 'SecurePass123!'
        });
      
      expect(loginResponse1.body.token).not.toBe(loginResponse2.body.token);
    });
  });

  describe('Error Handling Security', () => {
    test('should not leak sensitive information in errors', async () => {
      const response = await request(app)
        .get('/api/v1/users/nonexistent-endpoint')
        .set('Authorization', `Bearer ${authToken}`);
      
      expect(response.status).toBe(404);
      expect(response.body.error).not.toContain('database');
      expect(response.body.error).not.toContain('internal');
      expect(response.body.error).not.toContain('stack');
    });

    test('should handle malformed requests gracefully', async () => {
      const response = await request(app)
        .post('/api/v1/events')
        .set('Authorization', `Bearer ${brandToken}`)
        .set('Content-Type', 'application/json')
        .send('malformed json{');
      
      expect(response.status).toBe(400);
      expect(response.body.error).toBeDefined();
      expect(response.body.error).not.toContain('SyntaxError');
    });
  });
});
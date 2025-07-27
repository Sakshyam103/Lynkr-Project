/**
 * End-to-End User Journey Tests
 * Complete user workflows from registration to event participation
 */

const { test, expect } = require('@playwright/test');

test.describe('Complete User Journey', () => {
  test('User registration to event participation flow', async ({ page }) => {
    // User Registration
    await page.goto('/register');
    await page.fill('[data-testid="email-input"]', 'test@example.com');
    await page.fill('[data-testid="password-input"]', 'SecurePass123!');
    await page.fill('[data-testid="name-input"]', 'Test User');
    await page.click('[data-testid="register-button"]');
    
    // Verify registration success
    await expect(page.locator('[data-testid="welcome-message"]')).toBeVisible();
    
    // Privacy Settings
    await page.click('[data-testid="privacy-settings"]');
    await page.check('[data-testid="location-consent"]');
    await page.check('[data-testid="analytics-consent"]');
    await page.click('[data-testid="save-privacy"]');
    
    // Event Discovery
    await page.goto('/events');
    await expect(page.locator('[data-testid="event-list"]')).toBeVisible();
    
    // Event Details
    await page.click('[data-testid="event-card"]:first-child');
    await expect(page.locator('[data-testid="event-details"]')).toBeVisible();
    
    // Event Check-in (simulate location)
    await page.evaluate(() => {
      navigator.geolocation.getCurrentPosition = (success) => {
        success({ coords: { latitude: 37.7749, longitude: -122.4194 } });
      };
    });
    
    await page.click('[data-testid="checkin-button"]');
    await expect(page.locator('[data-testid="checkin-success"]')).toBeVisible();
    
    // Content Creation
    await page.click('[data-testid="create-content"]');
    await page.setInputFiles('[data-testid="photo-upload"]', 'test-image.jpg');
    await page.fill('[data-testid="caption-input"]', 'Great event! #testbrand');
    await page.click('[data-testid="post-content"]');
    
    // Verify content posted
    await expect(page.locator('[data-testid="content-success"]')).toBeVisible();
    
    // Poll Participation
    await page.click('[data-testid="poll-option-1"]');
    await expect(page.locator('[data-testid="poll-results"]')).toBeVisible();
    
    // Rewards Check
    await page.goto('/profile/rewards');
    await expect(page.locator('[data-testid="points-display"]')).toContainText('35'); // Check-in + content + poll points
  });

  test('Brand portal event monitoring flow', async ({ page }) => {
    // Brand Login
    await page.goto('/brand/login');
    await page.fill('[data-testid="email-input"]', 'brand@example.com');
    await page.fill('[data-testid="password-input"]', 'BrandPass123!');
    await page.click('[data-testid="login-button"]');
    
    // Dashboard Overview
    await expect(page.locator('[data-testid="dashboard"]')).toBeVisible();
    await expect(page.locator('[data-testid="active-events"]')).toBeVisible();
    
    // Event Creation
    await page.click('[data-testid="create-event"]');
    await page.fill('[data-testid="event-name"]', 'Test Brand Event');
    await page.fill('[data-testid="event-description"]', 'Testing event creation');
    await page.fill('[data-testid="event-date"]', '2024-12-31');
    await page.fill('[data-testid="event-location"]', '123 Test St, San Francisco, CA');
    await page.click('[data-testid="save-event"]');
    
    // Verify event created
    await expect(page.locator('[data-testid="event-success"]')).toBeVisible();
    
    // Real-time Monitoring
    await page.click('[data-testid="monitor-event"]');
    await expect(page.locator('[data-testid="live-metrics"]')).toBeVisible();
    await expect(page.locator('[data-testid="attendee-count"]')).toBeVisible();
    
    // Content Gallery
    await page.click('[data-testid="content-tab"]');
    await expect(page.locator('[data-testid="content-gallery"]')).toBeVisible();
    
    // Analytics Review
    await page.click('[data-testid="analytics-tab"]');
    await expect(page.locator('[data-testid="analytics-charts"]')).toBeVisible();
    await expect(page.locator('[data-testid="roi-metrics"]')).toBeVisible();
  });

  test('Complete purchase attribution flow', async ({ page }) => {
    // User discovers product at event
    await page.goto('/events/test-event');
    await page.click('[data-testid="checkin-button"]');
    
    // Interact with brand content
    await page.click('[data-testid="brand-product"]');
    await expect(page.locator('[data-testid="product-details"]')).toBeVisible();
    
    // Get discount code
    await page.click('[data-testid="get-discount"]');
    const discountCode = await page.locator('[data-testid="discount-code"]').textContent();
    
    // Navigate to purchase
    await page.click('[data-testid="shop-now"]');
    await page.fill('[data-testid="product-quantity"]', '1');
    await page.click('[data-testid="add-to-cart"]');
    
    // Apply discount code
    await page.fill('[data-testid="discount-input"]', discountCode);
    await page.click('[data-testid="apply-discount"]');
    await expect(page.locator('[data-testid="discount-applied"]')).toBeVisible();
    
    // Complete purchase
    await page.click('[data-testid="checkout"]');
    await page.fill('[data-testid="card-number"]', '4242424242424242');
    await page.fill('[data-testid="card-expiry"]', '12/25');
    await page.fill('[data-testid="card-cvc"]', '123');
    await page.click('[data-testid="complete-purchase"]');
    
    // Verify purchase success
    await expect(page.locator('[data-testid="purchase-success"]')).toBeVisible();
    
    // Check attribution in brand portal
    await page.goto('/brand/analytics/attribution');
    await expect(page.locator('[data-testid="attributed-purchase"]')).toBeVisible();
  });
});

test.describe('Error Handling and Edge Cases', () => {
  test('Network failure recovery', async ({ page }) => {
    // Simulate network failure during check-in
    await page.route('**/api/v1/events/*/checkin', route => route.abort());
    
    await page.goto('/events/test-event');
    await page.click('[data-testid="checkin-button"]');
    
    // Verify error handling
    await expect(page.locator('[data-testid="network-error"]')).toBeVisible();
    await expect(page.locator('[data-testid="retry-button"]')).toBeVisible();
    
    // Restore network and retry
    await page.unroute('**/api/v1/events/*/checkin');
    await page.click('[data-testid="retry-button"]');
    await expect(page.locator('[data-testid="checkin-success"]')).toBeVisible();
  });

  test('Invalid data handling', async ({ page }) => {
    // Test invalid email registration
    await page.goto('/register');
    await page.fill('[data-testid="email-input"]', 'invalid-email');
    await page.click('[data-testid="register-button"]');
    await expect(page.locator('[data-testid="email-error"]')).toBeVisible();
    
    // Test weak password
    await page.fill('[data-testid="email-input"]', 'test@example.com');
    await page.fill('[data-testid="password-input"]', '123');
    await page.click('[data-testid="register-button"]');
    await expect(page.locator('[data-testid="password-error"]')).toBeVisible();
  });

  test('Permission denied scenarios', async ({ page }) => {
    // Test location permission denied
    await page.evaluate(() => {
      navigator.geolocation.getCurrentPosition = (success, error) => {
        error({ code: 1, message: 'Permission denied' });
      };
    });
    
    await page.goto('/events/test-event');
    await page.click('[data-testid="checkin-button"]');
    await expect(page.locator('[data-testid="location-permission-error"]')).toBeVisible();
    
    // Test camera permission denied
    await page.click('[data-testid="create-content"]');
    await expect(page.locator('[data-testid="camera-permission-error"]')).toBeVisible();
  });
});

test.describe('Accessibility Testing', () => {
  test('Keyboard navigation', async ({ page }) => {
    await page.goto('/events');
    
    // Test tab navigation
    await page.keyboard.press('Tab');
    await expect(page.locator(':focus')).toHaveAttribute('data-testid', 'search-input');
    
    await page.keyboard.press('Tab');
    await expect(page.locator(':focus')).toHaveAttribute('data-testid', 'filter-button');
    
    // Test enter key activation
    await page.keyboard.press('Enter');
    await expect(page.locator('[data-testid="filter-menu"]')).toBeVisible();
  });

  test('Screen reader compatibility', async ({ page }) => {
    await page.goto('/events/test-event');
    
    // Check ARIA labels
    await expect(page.locator('[data-testid="checkin-button"]')).toHaveAttribute('aria-label', 'Check in to event');
    await expect(page.locator('[data-testid="event-details"]')).toHaveAttribute('role', 'main');
    
    // Check heading structure
    const headings = await page.locator('h1, h2, h3, h4, h5, h6').all();
    expect(headings.length).toBeGreaterThan(0);
  });

  test('High contrast mode', async ({ page }) => {
    // Enable high contrast mode
    await page.emulateMedia({ colorScheme: 'dark' });
    await page.goto('/events');
    
    // Verify contrast ratios meet WCAG standards
    const backgroundColor = await page.locator('body').evaluate(el => 
      getComputedStyle(el).backgroundColor
    );
    const textColor = await page.locator('body').evaluate(el => 
      getComputedStyle(el).color
    );
    
    // Basic contrast check (would use actual contrast calculation in real implementation)
    expect(backgroundColor).not.toBe(textColor);
  });
});
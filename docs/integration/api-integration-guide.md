# API Integration Guide

## Overview

The Lynkr API provides programmatic access to event data, user analytics, and content management features. This guide covers authentication, endpoints, and integration patterns.

## Authentication

### API Keys
```bash
# Include API key in header
curl -H "Authorization: Bearer YOUR_API_KEY" \
     -H "Content-Type: application/json" \
     https://api.lynkr.com/v1/events
```

### OAuth2 Flow
```javascript
// OAuth2 authentication flow
const authUrl = 'https://api.lynkr.com/oauth/authorize';
const params = {
  client_id: 'your_client_id',
  response_type: 'code',
  scope: 'events:read analytics:read',
  redirect_uri: 'https://yourapp.com/callback'
};
```

## Core Endpoints

### Events API
```javascript
// Get events
GET /v1/events
GET /v1/events/{eventId}

// Create event
POST /v1/events
{
  "name": "Product Launch Event",
  "description": "Launch of our new product line",
  "startDate": "2024-04-15T10:00:00Z",
  "endDate": "2024-04-15T18:00:00Z",
  "location": {
    "latitude": 37.7749,
    "longitude": -122.4194,
    "address": "123 Main St, San Francisco, CA"
  }
}
```

### Analytics API
```javascript
// Get event analytics
GET /v1/events/{eventId}/analytics/attendance
GET /v1/events/{eventId}/analytics/engagement
GET /v1/events/{eventId}/analytics/content

// Response example
{
  "eventId": "event_123",
  "totalAttendees": 1247,
  "uniqueUsers": 1180,
  "averageDuration": 3600,
  "contentCreated": 89,
  "engagementRate": 0.73
}
```

### Content API
```javascript
// Get event content
GET /v1/events/{eventId}/content

// Content permissions
PUT /v1/content/{contentId}/permissions
{
  "brandUsage": true,
  "publicDisplay": false,
  "attribution": "required"
}
```

## E-commerce Integration

### Shopify Integration
```javascript
// Shopify webhook setup
POST /v1/ecommerce/integrations
{
  "platform": "shopify",
  "storeUrl": "your-store.myshopify.com",
  "accessToken": "your_access_token",
  "webhookUrl": "https://api.lynkr.com/webhooks/shopify"
}

// Track purchase
POST /v1/ecommerce/purchases
{
  "orderId": "order_123",
  "userId": "user_456",
  "eventId": "event_789",
  "amount": 99.99,
  "currency": "USD",
  "products": [
    {
      "id": "product_123",
      "name": "Product Name",
      "price": 99.99,
      "quantity": 1
    }
  ]
}
```

### WooCommerce Integration
```php
// WooCommerce webhook handler
function lynkr_track_purchase($order_id) {
    $order = wc_get_order($order_id);
    
    $data = array(
        'orderId' => $order_id,
        'userId' => $order->get_customer_id(),
        'amount' => $order->get_total(),
        'currency' => $order->get_currency(),
        'eventId' => get_post_meta($order_id, '_lynkr_event_id', true)
    );
    
    wp_remote_post('https://api.lynkr.com/v1/ecommerce/purchases', array(
        'headers' => array(
            'Authorization' => 'Bearer ' . LYNKR_API_KEY,
            'Content-Type' => 'application/json'
        ),
        'body' => json_encode($data)
    ));
}
add_action('woocommerce_order_status_completed', 'lynkr_track_purchase');
```

## CRM Integration

### Salesforce Integration
```javascript
// Salesforce contact sync
POST /v1/crm/integrations
{
  "crmType": "salesforce",
  "instanceUrl": "https://yourorg.salesforce.com",
  "accessToken": "your_access_token",
  "refreshToken": "your_refresh_token"
}

// Sync event attendees
POST /v1/crm/salesforce/sync/{eventId}
{
  "objectType": "Contact",
  "fieldMapping": {
    "email": "Email",
    "name": "Name",
    "eventId": "Event_ID__c"
  }
}
```

### HubSpot Integration
```javascript
// HubSpot contact creation
const hubspotData = {
  properties: {
    email: attendee.email,
    firstname: attendee.firstName,
    lastname: attendee.lastName,
    event_attended: eventId,
    attendance_date: new Date().toISOString()
  }
};

fetch('https://api.hubapi.com/crm/v3/objects/contacts', {
  method: 'POST',
  headers: {
    'Authorization': `Bearer ${hubspotToken}`,
    'Content-Type': 'application/json'
  },
  body: JSON.stringify(hubspotData)
});
```

## Webhooks

### Setting Up Webhooks
```javascript
// Register webhook
POST /v1/webhooks
{
  "url": "https://yourapp.com/webhooks/lynkr",
  "events": ["event.checkin", "content.created", "purchase.completed"],
  "secret": "your_webhook_secret"
}
```

### Webhook Payload Examples
```javascript
// Event check-in webhook
{
  "event": "event.checkin",
  "timestamp": "2024-03-15T10:30:00Z",
  "data": {
    "eventId": "event_123",
    "userId": "user_456",
    "checkinTime": "2024-03-15T10:30:00Z",
    "location": {
      "latitude": 37.7749,
      "longitude": -122.4194
    }
  }
}

// Content created webhook
{
  "event": "content.created",
  "timestamp": "2024-03-15T11:15:00Z",
  "data": {
    "contentId": "content_789",
    "userId": "user_456",
    "eventId": "event_123",
    "mediaType": "photo",
    "tags": ["product_launch", "brand_name"]
  }
}
```

## SDK Examples

### JavaScript SDK
```javascript
import LynkrSDK from '@lynkr/sdk';

const lynkr = new LynkrSDK({
  apiKey: 'your_api_key',
  environment: 'production' // or 'sandbox'
});

// Track event
await lynkr.events.track({
  eventId: 'event_123',
  userId: 'user_456',
  action: 'checkin'
});

// Get analytics
const analytics = await lynkr.analytics.getEventMetrics('event_123');
```

### Python SDK
```python
from lynkr import LynkrClient

client = LynkrClient(api_key='your_api_key')

# Create event
event = client.events.create({
    'name': 'Product Launch',
    'start_date': '2024-04-15T10:00:00Z',
    'location': {
        'latitude': 37.7749,
        'longitude': -122.4194
    }
})

# Get analytics
analytics = client.analytics.get_event_metrics(event.id)
print(f"Total attendees: {analytics.total_attendees}")
```

## Rate Limits

- **Standard**: 1000 requests per hour
- **Premium**: 5000 requests per hour
- **Enterprise**: Custom limits available

## Error Handling

```javascript
// Error response format
{
  "error": {
    "code": "INVALID_REQUEST",
    "message": "The request is missing required parameters",
    "details": {
      "missing_fields": ["eventId", "userId"]
    }
  }
}

// Common error codes
// 400 - Bad Request
// 401 - Unauthorized
// 403 - Forbidden
// 404 - Not Found
// 429 - Rate Limited
// 500 - Internal Server Error
```

## Testing

### Sandbox Environment
- Base URL: `https://api-sandbox.lynkr.com/v1`
- Test API keys available in brand portal
- Sample data provided for testing

### Postman Collection
Download our Postman collection for easy API testing:
`https://api.lynkr.com/postman/collection.json`
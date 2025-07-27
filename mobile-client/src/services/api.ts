/**
 * API Service Layer - Added createEvent for brands
 */

const API_BASE_URL = 'http://localhost:8080';

class ApiService {
  private token: string | null = null;
  private userRole: string | null = null;

  setToken(token: string) {
    this.token = token;
  }

  setUserRole(role: string) {
    this.userRole = role;
  }

  private getBaseUrl() {
    switch (this.userRole) {
      case 'user':
        return `${API_BASE_URL}/user/v1`;
      case 'brand':
        return `${API_BASE_URL}/brand/v1`;
      default:
        return `${API_BASE_URL}/api/v1`;
    }
  }

  private async request(endpoint: string, options: RequestInit = {}) {
    const baseUrl = endpoint.startsWith('/api/v1') ? API_BASE_URL : this.getBaseUrl();
    const url = `${baseUrl}${endpoint}`;
    
    const headers = {
      'Content-Type': 'application/json',
      ...(this.token && { Authorization: `Bearer ${this.token}` }),
      ...options.headers,
    };

    const response = await fetch(url, { ...options, headers });
    
    if (!response.ok) {
      throw new Error(`API Error: ${response.status}`);
    }
    
    return response.json();
  }

  // FIXED: Events endpoint - call the actual backend endpoint
  async getEvents() {
    return this.request('/events', { method: 'GET' });
  }

  // ADDED: Create event method for brands
  async createEvent(eventData: any) {
    // This will call POST /brand/v1/events - you need to add this endpoint to your backend
    return this.request('/events', {
      method: 'POST',
      body: JSON.stringify(eventData),
    });
  }

  // Public Auth endpoints
  async register(userData: any) {
    return this.request('/api/v1/users/register', {
      method: 'POST',
      body: JSON.stringify(userData),
    });
  }

  async login(credentials: any) {
    return this.request('/api/v1/users/login', {
      method: 'POST',
      body: JSON.stringify(credentials),
    });
  }

  async brandLogin(credentials: any) {
    return this.request('/api/v1/brands/login', {
      method: 'POST',
      body: JSON.stringify(credentials),
    });
  }

  // USER ONLY ENDPOINTS
  async updatePrivacySettings(settings: any) {
    return this.request('/users/privacy', {
      method: 'PUT',
      body: JSON.stringify(settings),
    });
  }

  async checkIn(eventId: string, location: any) {
    return this.request(`/events/${eventId}/checkin`, {
      method: 'POST',
      body: JSON.stringify(location),
    });
  }

  async checkOut(eventId: string) {
    return this.request(`/events/${eventId}/checkout`, {
      method: 'POST',
    });
  }

  async getEventTags(eventId: string) {
    return this.request(`/events/${eventId}/tags`);
  }

  async createContent(contentData: any) {
    return this.request('/content', {
      method: 'POST',
      body: JSON.stringify(contentData),
    });
  }

  async uploadContent(formData: FormData) {
    const baseUrl = this.getBaseUrl();
    const url = `${baseUrl}/content`;
    
    const headers: any = {};
    if (this.token) {
      headers.Authorization = `Bearer ${this.token}`;
    }

    const response = await fetch(url, {
      method: 'POST',
      headers,
      body: formData,
    });
    
    if (!response.ok) {
      throw new Error(`API Error: ${response.status}`);
    }
    
    return response.json();
  }

  async updateContentPermissions(id: string, permissions: any) {
    return this.request(`user/v1/content/${id}/permissions`, {
      method: 'PUT',
      body: JSON.stringify(permissions),
    });
  }

  async getUserRewards() {
    return this.request('/users/rewards');
  }

  async getAvailableSurveys() {
    return this.request('/surveys/available');
  }

  async getUserBadges() {
    return this.request('/users/badges');
  }

  async submitPollVote(pollData: any) {
    return this.request('/feedback/polls/vote', {
      method: 'POST',
      body: JSON.stringify(pollData),
    });
  }

  async submitSliderFeedback(feedbackData: any) {
    return this.request('/feedback/sliders', {
      method: 'POST',
      body: JSON.stringify(feedbackData),
    });
  }

  async submitQuickFeedback(feedbackData: any) {
    return this.request('/feedback/quick', {
      method: 'POST',
      body: JSON.stringify(feedbackData),
    });
  }

  async trackPurchase(purchaseData: any) {
    return this.request('/ecommerce/purchases', {
      method: 'POST',
      body: JSON.stringify(purchaseData),
    });
  }

  async submitSurveyResponse(responseData: any) {
    return this.request('/surveys/respond', {
      method: 'POST',
      body: JSON.stringify(responseData),
    });
  }

  async validateDiscountCode(code: string) {
    return this.request(`/discount/codes/${code}/validate`);
  }

  async redeemDiscountCode(redeemData: any) {
    return this.request('/discount/redeem', {
      method: 'POST',
      body: JSON.stringify(redeemData),
    });
  }

  async getEvent(id: string) {
    return this.request(`/events/${id}`);
  }

  async getEventContent(eventId: string) {
    return this.request(`/events/${eventId}/content`);
  }

  async requestDataDeletion() {
    return this.request('/users/data/delete', {
      method: 'POST',
    });
  }

  async exportUserData() {
    return this.request('/users/data/export');
  }

  async anonymizeUserData() {
    return this.request('/users/data/anonymize', {
      method: 'POST',
    });
  }

  // BRAND ONLY ENDPOINTS
  async getBrandDashboard() {
    return this.request('/brands/dashboard');
  }

  async trackContentAnalytics(id: string, analytics: any) {
    return this.request(`/content/${id}/analytics`, {
      method: 'POST',
      body: JSON.stringify(analytics),
    });
  }

  async getContent(id: string) {
    return this.request(`/content/${id}`);
  }

  async searchTags(query: string) {
    return this.request(`/content/tags/search?q=${query}`);
  }

  async getBrandCampaigns() {
    return this.request('/brands/campaigns');
  }

  async createCampaign(campaignData: any) {
    return this.request('/brands/campaigns', {
      method: 'POST',
      body: JSON.stringify(campaignData),
    });
  }

  async getBrandContent() {
    return this.request('/brands/content');
  }

  async analyzeSentiment(text: string) {
    return this.request('/sentiment/analyze', {
      method: 'POST',
      body: JSON.stringify({ text }),
    });
  }

  async getEventSentiment(eventId: string) {
    return this.request(`/events/${eventId}/sentiment`);
  }

  async getEngagementMetrics(eventId: string) {
    return this.request(`/events/${eventId}/analytics/engagement`);
  }

  async getAttendanceAnalytics(eventId: string) {
    return this.request(`/events/${eventId}/analytics/attendance`);
  }

  async getContentPerformance(eventId: string) {
    return this.request(`/events/${eventId}/analytics/content`);
  }

  async getRealtimeStats(eventId: string) {
    return this.request(`/events/${eventId}/analytics/realtime`);
  }

  async trackEvent(eventData: any) {
    return this.request('/analytics/track', {
      method: 'POST',
      body: JSON.stringify(eventData),
    });
  }

  async createEcommerceIntegration(integrationData: any) {
    return this.request('/ecommerce/integrations', {
      method: 'POST',
      body: JSON.stringify(integrationData),
    });
  }

  async getEcommerceIntegration() {
    return this.request('/ecommerce/integrations');
  }

  async getPurchaseAnalytics(eventId: string) {
    return this.request(`/events/${eventId}/purchases/analytics`);
  }

  async getTopProducts(eventId: string) {
    return this.request(`/events/${eventId}/purchases/top-products`);
  }

  async generateDiscountCode(codeData: any) {
    return this.request('/discount/generate', {
      method: 'POST',
      body: JSON.stringify(codeData),
    });
  }

  async getCodeAnalytics(eventId: string) {
    return this.request(`/events/${eventId}/discount/analytics`);
  }

  async getBrandCodes() {
    return this.request('/brands/discount/codes');
  }

  async trackPixelSearch(searchData: any) {
    return this.request('/pixel/search', {
      method: 'POST',
      body: JSON.stringify(searchData),
    });
  }

  async getPixelAnalytics(eventId: string) {
    return this.request(`/events/${eventId}/pixel/analytics`);
  }

  async generatePixelURL() {
    return this.request('/pixel/generate');
  }

  async processContentAI(contentId: string) {
    return this.request(`/content/${contentId}/ai-process`, {
      method: 'POST',
    });
  }

  async getProductAnalytics() {
    return this.request('/brands/product-analytics');
  }

  async getConversionFunnel(eventId: string) {
    return this.request(`/events/${eventId}/conversion-funnel`);
  }

  async getAttributionReport(eventId: string) {
    return this.request(`/events/${eventId}/attribution-report`);
  }

  async awardReward(rewardData: any) {
    return this.request('/rewards/award', {
      method: 'POST',
      body: JSON.stringify(rewardData),
    });
  }

  async processQualityRewards(qualityData: any) {
    return this.request('/rewards/process-quality', {
      method: 'POST',
      body: JSON.stringify(qualityData),
    });
  }

  async scheduleSurveys(eventId: string, surveyData: any) {
    return this.request(`/events/${eventId}/surveys/schedule`, {
      method: 'POST',
      body: JSON.stringify(surveyData),
    });
  }

  async getSurveyAnalytics(eventId: string) {
    return this.request(`/events/${eventId}/surveys/analytics`);
  }

  async createExportRequest(exportData: any, datatype: string, format: string) {
    const requestBody = {
      exportData,
      datatype,
      format
    };
    return this.request('/export/create', {
      method: 'POST',
      body: JSON.stringify(requestBody),
    });
  }

  async getExportStatus(requestId: string) {
    return this.request(`/export/${requestId}/status`);
  }

  async getExportFormats() {
    return this.request('/export/formats');
  }

  async createCRMIntegration(crmData: any) {
    return this.request('/crm/integrations', {
      method: 'POST',
      body: JSON.stringify(crmData),
    });
  }

  async syncEventData(integrationId: string, eventId: string) {
    return this.request(`/crm/${integrationId}/sync/${eventId}`, {
      method: 'POST',
    });
  }

  async getCRMTypes() {
    return this.request('/crm/types');
  }

  // PUBLIC ENDPOINTS
  async trackConversion(conversionData: any) {
    return this.request('/api/v1/conversions/track', {
      method: 'POST',
      body: JSON.stringify(conversionData),
    });
  }

  // Logout method
  logout() {
    this.token = null;
    this.userRole = null;
  }
}

export const apiService = new ApiService();
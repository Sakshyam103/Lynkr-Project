import React, { useState } from 'react';
import { View, Text, ScrollView, TouchableOpacity, StyleSheet, Alert } from 'react-native';
import { useSelector } from 'react-redux';
import { RootState } from '../store/store';
import { apiService } from '../services/api';

export default function EndpointTestScreen() {
  const { user } = useSelector((state: RootState) => state.auth);
  const [testResults, setTestResults] = useState<any>({});
  const [testing, setTesting] = useState(false);

  const isUser = user?.role === 'user' || !user?.role;
  const isBrand = user?.role === 'brand';

  const testEndpoint = async (name: string, testFunction: () => Promise<any>) => {
    try {
      setTesting(true);
      const result = await testFunction();
      setTestResults(prev => ({ ...prev, [name]: { success: true, data: result } }));
    } catch (error) {
      setTestResults(prev => ({ ...prev, [name]: { success: false, error: error.message } }));
    } finally {
      setTesting(false);
    }
  };

  const userEndpoints = [
    {
      name: 'Check In Event',
      test: () => apiService.checkIn('1', { latitude: 37.7749, longitude: -122.4194 })
    },
    {
      name: 'Get Event Tags',
      test: () => apiService.getEventTags('1')
    },
    {
      name: 'Create Content',
      test: () => apiService.createContent({
        eventId: '1',
        mediaType: 'photo',
        caption: 'Test content',
        tags: ['test']
      })
    },
    {
      name: 'Get User Rewards',
      test: () => apiService.getUserRewards()
    },
    {
      name: 'Get Available Surveys',
      test: () => apiService.getAvailableSurveys()
    },
    {
      name: 'Get User Badges',
      test: () => apiService.getUserBadges()
    },
    {
      name: 'Submit Poll Vote',
      test: () => apiService.submitPollVote({
        pollId: 'poll_1',
        optionId: 1,
        eventId: '1'
      })
    },
    {
      name: 'Submit Slider Feedback',
      test: () => apiService.submitSliderFeedback({
        eventId: '1',
        category: 'overall',
        rating: 4.5
      })
    },
    {
      name: 'Track Purchase',
      test: () => apiService.trackPurchase({
        orderId: 'order_123',
        eventId: '1',
        amount: 99.99
      })
    },
    {
      name: 'Validate Discount Code',
      test: () => apiService.validateDiscountCode('SAVE20')
    },
    {
      name: 'Get Event',
      test: () => apiService.getEvent('1')
    },
    {
      name: 'Get Event Content',
      test: () => apiService.getEventContent('1')
    },
    {
      name: 'Export User Data',
      test: () => apiService.exportUserData()
    }
  ];

  const brandEndpoints = [
    {
      name: 'Get Brand Dashboard',
      test: () => apiService.getBrandDashboard()
    },
    {
      name: 'Get Brand Campaigns',
      test: () => apiService.getBrandCampaigns()
    },
    {
      name: 'Create Campaign',
      test: () => apiService.createCampaign({
        name: 'Test Campaign',
        description: 'Test campaign description',
        budget: 10000
      })
    },
    {
      name: 'Get Brand Content',
      test: () => apiService.getBrandContent()
    },
    {
      name: 'Analyze Sentiment',
      test: () => apiService.analyzeSentiment('This is amazing!')
    },
    {
      name: 'Get Event Sentiment',
      test: () => apiService.getEventSentiment('1')
    },
    {
      name: 'Get Engagement Metrics',
      test: () => apiService.getEngagementMetrics('1')
    },
    {
      name: 'Get Attendance Analytics',
      test: () => apiService.getAttendanceAnalytics('1')
    },
    {
      name: 'Get Realtime Stats',
      test: () => apiService.getRealtimeStats('1')
    },
    {
      name: 'Generate Discount Code',
      test: () => apiService.generateDiscountCode({
        eventId: '1',
        discountType: 'percentage',
        discountValue: 20,
        maxUses: 100
      })
    },
    {
      name: 'Get Brand Codes',
      test: () => apiService.getBrandCodes()
    },
    {
      name: 'Get Product Analytics',
      test: () => apiService.getProductAnalytics()
    },
    {
      name: 'Get Conversion Funnel',
      test: () => apiService.getConversionFunnel('1')
    },
    {
      name: 'Create Export Request',
      test: () => apiService.createExportRequest({
        eventId: '1',
        format: 'CSV',
        fields: ['email', 'name']
      })
    }
  ];

  const testAllEndpoints = async () => {
    const endpoints = isUser ? userEndpoints : brandEndpoints;
    
    for (const endpoint of endpoints) {
      await testEndpoint(endpoint.name, endpoint.test);
      // Small delay between tests
      await new Promise(resolve => setTimeout(resolve, 500));
    }
    
    Alert.alert('Testing Complete', 'All endpoint tests have been completed. Check results below.');
  };

  const getResultColor = (result: any) => {
    if (!result) return '#999';
    return result.success ? '#28a745' : '#dc3545';
  };

  const getResultIcon = (result: any) => {
    if (!result) return '⏳';
    return result.success ? '✅' : '❌';
  };

  const endpoints = isUser ? userEndpoints : brandEndpoints;

  return (
    <ScrollView style={styles.container}>
      <View style={styles.header}>
        <Text style={styles.title}>Endpoint Testing</Text>
        <Text style={styles.subtitle}>
          Testing {isUser ? 'User' : 'Brand'} endpoints for role: {user?.role?.toUpperCase()}
        </Text>
        
        <TouchableOpacity 
          style={[styles.testAllButton, testing && styles.disabled]}
          onPress={testAllEndpoints}
          disabled={testing}
        >
          <Text style={styles.testAllText}>
            {testing ? 'Testing...' : 'Test All Endpoints'}
          </Text>
        </TouchableOpacity>
      </View>

      <View style={styles.endpointsList}>
        {endpoints.map((endpoint, index) => {
          const result = testResults[endpoint.name];
          return (
            <View key={index} style={styles.endpointCard}>
              <View style={styles.endpointHeader}>
                <Text style={styles.endpointName}>{endpoint.name}</Text>
                <Text style={styles.endpointStatus}>
                  {getResultIcon(result)}
                </Text>
              </View>
              
              {result && (
                <View style={styles.resultContainer}>
                  <Text style={[styles.resultStatus, { color: getResultColor(result) }]}>
                    {result.success ? 'SUCCESS' : 'FAILED'}
                  </Text>
                  {result.error && (
                    <Text style={styles.errorText}>{result.error}</Text>
                  )}
                  {result.success && result.data && (
                    <Text style={styles.successText}>
                      Response received: {typeof result.data === 'object' ? 'Object' : result.data}
                    </Text>
                  )}
                </View>
              )}
              
              <TouchableOpacity 
                style={styles.testButton}
                onPress={() => testEndpoint(endpoint.name, endpoint.test)}
                disabled={testing}
              >
                <Text style={styles.testButtonText}>Test Individual</Text>
              </TouchableOpacity>
            </View>
          );
        })}
      </View>

      <View style={styles.summary}>
        <Text style={styles.summaryTitle}>Test Summary</Text>
        <Text style={styles.summaryText}>
          Total Endpoints: {endpoints.length}
        </Text>
        <Text style={styles.summaryText}>
          Tested: {Object.keys(testResults).length}
        </Text>
        <Text style={styles.summaryText}>
          Passed: {Object.values(testResults).filter((r: any) => r.success).length}
        </Text>
        <Text style={styles.summaryText}>
          Failed: {Object.values(testResults).filter((r: any) => !r.success).length}
        </Text>
      </View>
    </ScrollView>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#f8f9fa',
  },
  header: {
    padding: 20,
    backgroundColor: '#fff',
  },
  title: {
    fontSize: 24,
    fontWeight: 'bold',
    color: '#333',
  },
  subtitle: {
    fontSize: 16,
    color: '#666',
    marginTop: 4,
    marginBottom: 20,
  },
  testAllButton: {
    backgroundColor: '#007AFF',
    padding: 16,
    borderRadius: 8,
    alignItems: 'center',
  },
  testAllText: {
    color: '#fff',
    fontSize: 16,
    fontWeight: '600',
  },
  disabled: {
    opacity: 0.6,
  },
  endpointsList: {
    padding: 20,
  },
  endpointCard: {
    backgroundColor: '#fff',
    borderRadius: 12,
    padding: 16,
    marginBottom: 12,
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.1,
    shadowRadius: 4,
    elevation: 3,
  },
  endpointHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: 12,
  },
  endpointName: {
    fontSize: 16,
    fontWeight: '600',
    color: '#333',
    flex: 1,
  },
  endpointStatus: {
    fontSize: 20,
  },
  resultContainer: {
    backgroundColor: '#f8f9fa',
    padding: 12,
    borderRadius: 8,
    marginBottom: 12,
  },
  resultStatus: {
    fontSize: 14,
    fontWeight: 'bold',
    marginBottom: 4,
  },
  errorText: {
    fontSize: 12,
    color: '#dc3545',
    fontFamily: 'monospace',
  },
  successText: {
    fontSize: 12,
    color: '#28a745',
  },
  testButton: {
    backgroundColor: '#28a745',
    padding: 8,
    borderRadius: 6,
    alignItems: 'center',
  },
  testButtonText: {
    color: '#fff',
    fontSize: 14,
    fontWeight: '600',
  },
  summary: {
    backgroundColor: '#fff',
    margin: 20,
    padding: 16,
    borderRadius: 12,
  },
  summaryTitle: {
    fontSize: 18,
    fontWeight: 'bold',
    marginBottom: 12,
  },
  summaryText: {
    fontSize: 14,
    color: '#666',
    marginBottom: 4,
  },
});